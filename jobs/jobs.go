package jobs

import (
	"context"
	"sync"
	"time"
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type JobStruct struct {
	JobsRepository models.JobsSvcRepo
}

func NewJobService() JobSvc {
	return &JobStruct{
		JobsRepository: models.JobsRepository(connections.PgDBConnection.Client),
	}
}

type JobSvc interface {
	FirstTimeJob(ctx context.Context, args []string)
	RetryJobs(ctx context.Context, args []string)
	JobConsumer(wg *sync.WaitGroup, ctx context.Context, jobsChannel chan models.ModelJobs)
	DequeueJobs(jobsChannel chan models.ModelJobs, status string)
}

func (j *JobStruct) JobConsumer(wg *sync.WaitGroup, ctx context.Context, jobsChannel chan models.ModelJobs) {
	defer wg.Done()
	serverTime := time.Now()
	for job := range jobsChannel {
		job.Status = constants.ProcessingJobStatus
		if err := j.JobsRepository.BulkUpsert([]*models.ModelJobs{&job}); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert jobs")
			continue
		}

		var jobError error
		switch job.JobType {
		case constants.InsertCsvFile:
			if err := ProcessInsertCsvFile(&job); err != nil {
				jobError = err
			}
		case constants.ExportCsvFile:
			if err := ProcessExportCsvFile(&job); err != nil {
				jobError = err
			}
		default:
			jobError = constants.InvalidJobTypeError(job.JobType)
		}

		if jobError != nil {
			retryAfter := serverTime.Add(time.Duration(job.RetryInterval) * time.Second)
			job.Status = constants.FailedJobStatus
			job.AddRuntimeError(jobError.Error())
			job.RunAfter = &retryAfter
		} else {
			job.Status = constants.CompletedJobStatus
		}

		if err := j.JobsRepository.BulkUpsert([]*models.ModelJobs{&job}); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert jobs")
		}
	}
}

func (j *JobStruct) DequeueJobs(jobsChannel chan models.ModelJobs, status string) {
	jobs := make([]*models.ModelJobs, 0, len(jobsChannel))
	for job := range jobsChannel {
		job.Status = status
		jobs = append(jobs, &job)
	}
	if len(jobs) > 0 {
		if err := j.JobsRepository.BulkUpsert(jobs); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert jobs")
		}
	}
	log.Info().Msgf("Dequeued %d jobs with status: %s", len(jobs), status)
}

func (j *JobStruct) FirstTimeJob(ctx context.Context, args []string) {
	var wg sync.WaitGroup
	jobsChannel := make(chan models.ModelJobs, 1000)

	ticker := time.NewTicker(time.Duration(conf.JobConfig.TickerInterval) * time.Second)
	defer func() {
		ticker.Stop()
		wg.Wait()
		log.Info().Msg("All workers stopped")
	}()

	for i := 0; i < conf.JobConfig.ParallelJobs; i++ {
		wg.Add(1)
		go j.JobConsumer(&wg, ctx, jobsChannel)
	}

	inQueSize := conf.JobConfig.JobInQueuedSize
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, stopping job producer...")
			close(jobsChannel)
			j.DequeueJobs(jobsChannel, constants.OpenJobStatus)
			return
		case <-ticker.C:
			if len(jobsChannel) >= inQueSize {
				continue
			}
			jobs, err := j.JobsRepository.ListByFilters(models.JobsFilters{
				Status: []string{constants.OpenJobStatus},
				Limit:  1,
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
			if err = j.JobsRepository.BulkUpsert(jobs); err != nil {
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

func (j *JobStruct) RetryJobs(ctx context.Context, args []string) {
	var wg sync.WaitGroup
	jobsChannel := make(chan models.ModelJobs, 1000)

	wg.Add(1)
	go j.JobConsumer(&wg, ctx, jobsChannel)

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
			j.DequeueJobs(jobsChannel, constants.FailedJobStatus)
			return
		case <-ticker.C:
			jobs, err := j.JobsRepository.ListByFilters(models.JobsFilters{
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
			if err = j.JobsRepository.BulkUpsert(jobs); err != nil {
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

func RunJobs(ctx context.Context, args []string) {
	if len(args) == 0 {
		log.Error().Msg("No job type provided")
		return
	}
	jobService := NewJobService()
	jobType := args[0]

	args = args[1:]
	switch jobType {
	case constants.FirstTimeJobType:
		jobService.FirstTimeJob(ctx, args)
	case constants.RetryJobType:
		jobService.RetryJobs(ctx, args)
	default:
		log.Error().Msgf("Invalid job type: %s", jobType)
	}
}
