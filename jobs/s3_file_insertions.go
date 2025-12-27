package jobs

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	commonService "vivek-ray/modules/common/service"
	"vivek-ray/utilities"

	"github.com/rs/zerolog/log"
)

type InsertJobStruct struct {
	JobsRepository     models.JobsSvcRepo
	BatchUpsertService commonService.BatchUpsertSvc
}

func NewInsertJob() *InsertJobStruct {
	return &InsertJobStruct{
		JobsRepository:     models.JobsRepository(connections.PgDBConnection.Client),
		BatchUpsertService: commonService.NewBatchUpsertService(),
	}
}

func rowToMap(headers, row []string) map[string]string {
	rowMap := make(map[string]string, len(headers))
	for idx, header := range headers {
		if idx < len(row) {
			rowMap[header] = row[idx]
		}
	}
	return rowMap
}

func (i *InsertJobStruct) Run(wg *sync.WaitGroup, retry int, jobsChannel chan models.ModelJobs) {
	defer wg.Done()

	for job := range jobsChannel {
		err := i.processJob(&job)

		if err != nil {
			job.Status = constants.FailedJobStatus
			job.RuntimeErrors = append(job.RuntimeErrors, err.Error())
			job.RetryCount = job.RetryCount - retry
		} else {
			job.Status = constants.CompletedJobStatus
		}

		if err := i.JobsRepository.BulkUpsert([]*models.ModelJobs{&job}); err != nil {
			log.Error().Err(err).Msg("Failed to update job status")
		}
	}
}

func (i *InsertJobStruct) processJob(job *models.ModelJobs) error {
	var jobData utilities.InsertFileJobData
	if err := json.Unmarshal(job.Data, &jobData); err != nil {
		return err
	}
	job.Status = constants.ProcessingJobStatus
	if err := i.JobsRepository.BulkUpsert([]*models.ModelJobs{job}); err != nil {
		return err
	}
	if jobData.FileS3Bucket == "" { // take default bucket from config if not provided
		jobData.FileS3Bucket = conf.S3StorageConfig.S3Bucket
	}
	fileStream, err := connections.S3Connection.ReadFileStream(
		context.Background(),
		jobData.FileS3Bucket,
		jobData.FileS3Key,
	)
	if err != nil {
		return err
	}
	defer fileStream.Close()

	switch jobData.FileType {
	case constants.FileTypeCsv:
		return i.processCSV(fileStream)
	}

	return nil
}

func (i *InsertJobStruct) processCSV(reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	headers, err := csvReader.Read()
	if err != nil {
		return err
	}
	log.Info().Msgf("Headers: %v", headers)
	batchSize := conf.JobConfig.BatchSize
	batch := make([]map[string]string, 0, batchSize)

	for {
		row, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			log.Info().Msg("EOF")
			break
		}
		if err != nil {
			return err
		}

		batch = append(batch, rowToMap(headers, row))
		if len(batch) >= batchSize {
			if err := i.BatchUpsertService.ProcessBatchUpsert(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		return i.BatchUpsertService.ProcessBatchUpsert(batch)
	}
	return nil
}

func dequeueJobs(jobsChannel chan models.ModelJobs, status string, insertJob *InsertJobStruct) {
	jobs := make([]*models.ModelJobs, 0, len(jobsChannel))

	for job := range jobsChannel {
		job.Status = status
		jobs = append(jobs, &job)
	}

	if len(jobs) > 0 {
		if err := insertJob.JobsRepository.BulkUpsert(jobs); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert jobs")
		}
		log.Info().Msgf("Dequeued %d jobs with status: %s", len(jobs), status)
	}
}

func InsertFileJob(ctx context.Context) {
	insertJob := NewInsertJob()
	var wg sync.WaitGroup
	jobsChannel := make(chan models.ModelJobs, 1000)

	ticker := time.NewTicker(time.Duration(conf.JobConfig.TickerInterval) * time.Minute)
	defer func() {
		ticker.Stop()
		wg.Wait()
		log.Info().Msg("All workers stopped")
	}()

	for i := 0; i < conf.JobConfig.ParallelJobs; i++ {
		wg.Add(1)
		go insertJob.Run(&wg, 0, jobsChannel)
	}

	inQueSize := conf.JobConfig.JobInQueuedSize
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, stopping job producer...")
			close(jobsChannel)
			dequeueJobs(jobsChannel, constants.OpenJobStatus, insertJob)
			return
		case <-ticker.C:
			if len(jobsChannel) >= inQueSize {
				continue
			}
			jobs, err := insertJob.JobsRepository.ListByFilters(models.JobsFilters{
				JobType: constants.InsertFileJobType,
				Status:  []string{constants.OpenJobStatus},
				Limit:   1,
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to list jobs")
				continue
			}
			if len(jobs) == 0 {
				log.Info().Msg("No jobs found to insert file")
				continue
			}

			for _, job := range jobs {
				log.Info().Msgf("pushing job to in queue: %s", job.UUID)
				job.Status = constants.InQueueJobStatus
			}
			if err = insertJob.JobsRepository.BulkUpsert(jobs); err != nil {
				log.Error().Err(err).Msg("Failed to bulk upsert jobs")
				continue
			}

			for _, job := range jobs {
				jobsChannel <- *job
				log.Info().Msgf("Job pushed to channel: %s", job.UUID)
			}

		}
	}
}

func RetryInsertFileJob(ctx context.Context) {
	insertJob := NewInsertJob()
	var wg sync.WaitGroup
	jobsChannel := make(chan models.ModelJobs, 1000)

	wg.Add(1)
	go insertJob.Run(&wg, 1, jobsChannel)

	ticker := time.NewTicker(time.Duration(conf.JobConfig.TickerInterval) * time.Minute)
	defer func() {
		ticker.Stop()
		wg.Wait()
		log.Info().Msg("All workers stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, stopping retry job producer...")
			close(jobsChannel)
			dequeueJobs(jobsChannel, constants.FailedJobStatus, insertJob)
			return
		case <-ticker.C:
			jobs, err := insertJob.JobsRepository.ListByFilters(models.JobsFilters{
				JobType:  constants.InsertFileJobType,
				Retrying: true,
				Status:   []string{constants.FailedJobStatus},
				Limit:    1,
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to list jobs")
				continue
			}
			if len(jobs) == 0 {
				log.Info().Msg("No jobs found to insert file")
				continue
			}

			for _, job := range jobs {
				job.Status = constants.RetryInQueuedJobStatus
			}
			if err = insertJob.JobsRepository.BulkUpsert(jobs); err != nil {
				log.Error().Err(err).Msg("Failed to bulk upsert jobs")
				continue
			}

			for _, job := range jobs {
				jobsChannel <- *job
				log.Info().Msgf("Job pushed to channel: %s", job.UUID)
			}
		}
	}
}
