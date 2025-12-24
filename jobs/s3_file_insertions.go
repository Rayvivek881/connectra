package jobs

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"sync"
	"time"
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/utilities"

	"github.com/rs/zerolog/log"
)

type InsertJobStruct struct {
	JobsRepository           models.JobsSvcRepo
	CompanyRepository        models.PgCompanySvcRepo
	ContactRepository        models.PgContactSvcRepo
	CompanyElasticRepository models.ElasticCompanySvcRepo
	ContactElasticRepository models.ElasticContactSvcRepo
}

func NewInsertJob() *InsertJobStruct {
	return &InsertJobStruct{
		JobsRepository:           models.JobsRepository(connections.PgDBConnection.Client),
		CompanyRepository:        models.PgCompanyRepository(connections.PgDBConnection.Client),
		ContactRepository:        models.PgContactRepository(connections.PgDBConnection.Client),
		CompanyElasticRepository: models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client),
		ContactElasticRepository: models.ElasticContactRepository(connections.ElasticsearchConnection.Client),
	}
}

func (i *InsertJobStruct) ProcessBatchInsert(batch []map[string]string) error {
	PgCompany := make([]*models.PgCompany, 0)
	PgContact := make([]*models.PgContact, 0)
	ElasticCompany := make([]*models.ElasticCompany, 0)
	ElasticContact := make([]*models.ElasticContact, 0)

	for _, row := range batch {
		company := models.PgCompanyFromRawData(row)
		contact := models.PgContactFromRowData(row, company)
		elasticCompany := models.ElasticCompanyFromRawData(company)
		elasticContact := models.ElasticContactFromRawData(contact, company)

		PgCompany = append(PgCompany, company)
		PgContact = append(PgContact, contact)
		ElasticCompany = append(ElasticCompany, elasticCompany)
		ElasticContact = append(ElasticContact, elasticContact)

	}
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		if _, err := i.CompanyRepository.BulkUpsert(PgCompany); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert companies")
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := i.ContactRepository.BulkUpsert(PgContact); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert contacts")
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := i.CompanyElasticRepository.BulkUpsert(ElasticCompany); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert companies to elasticsearch")
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := i.ContactElasticRepository.BulkUpsert(ElasticContact); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert contacts to elasticsearch")
		}
	}()
	wg.Wait()
	return nil
}

func (i *InsertJobStruct) Run(wg *sync.WaitGroup, jobsChannel chan models.ModelJobs) {
	defer wg.Done()

	for job := range jobsChannel {

		err := func() error {
			var Jobdata utilities.InsertFileJobData
			if err := json.Unmarshal(job.Data, &Jobdata); err != nil {
				return err
			}

			job.Status = constants.ProcessingJobStatus
			if err := i.JobsRepository.BulkUpsert([]*models.ModelJobs{&job}); err != nil {
				return err
			}

			batchSize, batch := 500, make([]map[string]string, 0)
			log.Info().Msgf("Processing job: %s", job.UUID)

			fileStream, err := connections.S3Connection.ReadFileStream(context.Background(), conf.S3StorageConfig.S3Bucket, Jobdata.FileS3Key)
			if err != nil {
				return err
			}
			defer fileStream.Close()
			switch Jobdata.FileType {
			case constants.FileTypeCsv:
				csvReader := csv.NewReader(fileStream)
				headers, err := csvReader.Read()
				if err != nil {
					return err
				}
				for {
					row, err := csvReader.Read()
					if err == io.EOF {
						break
					}

					if err != nil {
						return err
					}
					rowMap := make(map[string]string)
					for idx, header := range headers {
						if idx >= len(row) {
							continue
						}
						rowMap[utilities.GetCleanedString(header)] = utilities.GetCleanedString(row[idx])
					}
					batch = append(batch, rowMap)
					if len(batch) >= batchSize {
						if err := i.ProcessBatchInsert(batch); err != nil {
							return err
						}
						batch = make([]map[string]string, 0)
					}
				}
				if len(batch) > 0 {
					if err := i.ProcessBatchInsert(batch); err != nil {
						return err
					}
				}
			}
			return nil
		}()

		if err != nil {
			job.Status = constants.FailedJobStatus
			job.RuntimeErrors = append(job.RuntimeErrors, err.Error())
		} else {
			job.Status = constants.CompletedJobStatus
		}

		if err := i.JobsRepository.BulkUpsert([]*models.ModelJobs{&job}); err != nil {
			log.Error().Err(err).Msg("Failed to update job status")
		}
	}
}

func InsertFileJob() {
	insertJob := NewInsertJob()
	var wg sync.WaitGroup
	jobsChannel := make(chan models.ModelJobs, 1000)

	ticker := time.NewTicker(1 * time.Minute)
	defer func() {
		ticker.Stop()
		close(jobsChannel)
		wg.Wait()
		log.Info().Msg("All workers stopped")
	}()

	for i := 0; i < conf.JobConfig.ParallelJobs; i++ {
		wg.Add(1)
		go insertJob.Run(&wg, jobsChannel)
	}

	inQueSize := conf.JobConfig.JobInQueuedSize
	for range ticker.C {
		if len(jobsChannel) >= inQueSize {
			continue
		}

		jobs, err := insertJob.JobsRepository.ListByFilters(models.JobsFilters{
			JobType: constants.InsertFileJobType,
			Status:  []string{constants.OpenJobStatus},
			Limit:   inQueSize - len(jobsChannel),
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
			job.Status = constants.InQueueJobStatus
			jobsChannel <- *job
			log.Info().Msgf("Job pushed to channel: %s", job.UUID)
		}

		if err = insertJob.JobsRepository.BulkUpsert(jobs); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert jobs")
		}
	}
}
