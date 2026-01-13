package service

import (
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/modules/common/helper"

	"github.com/google/uuid"
)

type JobSvc interface {
	CreateJob(request helper.CreateJobRequest) error
	ListJobs(request helper.ListJobsRequest) ([]*models.ModelJobs, error)
	GetJobByUuid(jobUuid string) (*models.ModelJobs, error)
}

type jobService struct {
	jobsRepository models.JobsSvcRepo
}

func NewJobService() JobSvc {
	return &jobService{
		jobsRepository: models.JobsRepository(connections.PgDBConnection.Client),
	}
}

func (s *jobService) CreateJob(request helper.CreateJobRequest) error {
	job := &models.ModelJobs{
		UUID:       uuid.New().String(),
		JobType:    request.JobType,
		Data:       request.JobData,
		RetryCount: request.RetryCount,
	}
	return s.jobsRepository.BulkUpsert([]*models.ModelJobs{job})
}

func (s *jobService) ListJobs(request helper.ListJobsRequest) ([]*models.ModelJobs, error) {
	return s.jobsRepository.ListByFilters(models.JobsFilters{
		JobType: request.JobType,
		Status:  request.Status,
		Limit:   request.Limit,
	})
}

func (s *jobService) GetJobByUuid(jobUuid string) (*models.ModelJobs, error) {
	jobs, err := s.jobsRepository.ListByFilters(models.JobsFilters{
		Uuids: []string{jobUuid},
	})
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return nil, constants.JobNotFoundError
	}
	return jobs[0], nil
}
