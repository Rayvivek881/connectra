package jobs

import (
	"sync"
	"time"
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type S3FileInsertJob struct {
	companyPgRepository      models.PgCompanySvcRepo
	companyElasticRepository models.ElasticCompanySvcRepo
	contactPgRepository      models.PgContactSvcRepo
	contactElasticRepository models.ElasticContactSvcRepo
}

func NewS3FileInsertJob() *S3FileInsertJob {
	companyPgRepository := models.PgCompanyRepository(connections.PgDBConnection.Client)
	companyElasticRepository := models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client)
	contactPgRepository := models.PgContactRepository(connections.PgDBConnection.Client)
	contactElasticRepository := models.ElasticContactRepository(connections.ElasticsearchConnection.Client)
	return &S3FileInsertJob{
		companyPgRepository:      companyPgRepository,
		companyElasticRepository: companyElasticRepository,
		contactPgRepository:      contactPgRepository,
		contactElasticRepository: contactElasticRepository,
	}
}
func (j *S3FileInsertJob) Execute(job *models.ModelJobsData) error {
	return nil
}

func (j *S3FileInsertJob) Worker(wg *sync.WaitGroup, jobChannel chan *models.ModelJobsData) {
	wg.Add(1)
	defer wg.Done()

}

func (j *S3FileInsertJob) RetryWorker(wg *sync.WaitGroup, retryChannel chan *models.ModelJobsData) {
	wg.Add(1)
	defer wg.Done()

}

func InitS3FileInsertJob() {
	job := NewS3FileInsertJob()
	var jobWg sync.WaitGroup

	jobChannel := make(chan *models.ModelJobsData, conf.AppConfig.InQueueSize)
	for i := 0; i < conf.AppConfig.ParallelJobs; i++ {
		go job.Worker(&jobWg, jobChannel)
	}

	var retryWg sync.WaitGroup
	retryChannel := make(chan *models.ModelJobsData, conf.AppConfig.InQueueSize)
	go job.RetryWorker(&retryWg, retryChannel)

	for t := range time.Tick(time.Duration(conf.AppConfig.TickerInterval) * time.Minute) {
		log.Info().Msgf("Ticker interval: %s", t.Format("2006-01-02 15:04:05"))

	}

	defer func() {
		log.Info().Msgf("Waiting for jobs to complete")
		jobWg.Wait()
		retryWg.Wait()
		log.Info().Msgf("Jobs completed")
		close(jobChannel)
		close(retryChannel)
	}()
}
