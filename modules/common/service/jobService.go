package service

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/modules/common/helper"

	"github.com/uptrace/bun"
)

type JobService struct {
	jobsRepo  models.JobNodeSvcRepo
	edgesRepo models.EdgesSvcRepo
}

func NewJobService(db *bun.DB) JobServiceRepo {
	return &JobService{
		jobsRepo:  models.JobNodeRepository(db),
		edgesRepo: models.EdgesRepository(db),
	}
}

type JobServiceRepo interface {
	BulkInsertJobs(nodes []*helper.BulkInsertGraphRequest) error
	GetJobsCount(filters models.JobFilters) (int, error)
	GetJobs(filters *models.JobFilters) ([]*models.ModelJobNodes, error)
	InsertJobsToDB(jobs []*models.ModelJobNodes, edges []*models.ModelEdges) error
	UpdateAndRetriggerJob(uuid string, data json.RawMessage, retryCount *int) error
	GetJobByUUID(uuid string) (*models.ModelJobNodes, error)
}

func (j *JobService) GetJobsCount(filters models.JobFilters) (int, error) {
	return j.jobsRepo.GetJobsCount(&filters)
}

func (j *JobService) GetJobs(filters *models.JobFilters) ([]*models.ModelJobNodes, error) {
	return j.jobsRepo.GetJobs(filters)
}

func (j *JobService) InsertJobsToDB(jobs []*models.ModelJobNodes, edges []*models.ModelEdges) error {
	var responseErrors error
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := j.jobsRepo.JobsBulkUpsert(jobs)
		if err != nil {
			mu.Lock()
			responseErrors = errors.Join(responseErrors, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		err := j.edgesRepo.CreateEdges(edges)
		if err != nil {
			mu.Lock()
			responseErrors = errors.Join(responseErrors, err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	return responseErrors
}

func (j *JobService) BulkInsertJobs(jobNodes []*helper.BulkInsertGraphRequest) error {
	jobUuids := make([]string, len(jobNodes))
	for i, node := range jobNodes {
		jobUuids[i] = node.UUID
	}
	count, err := j.GetJobsCount(models.JobFilters{Uuids: jobUuids})
	if err != nil {
		return err
	}
	if count != 0 {
		return constants.ErrDuplicateJobUUIDs
	}
	jobs, edges := make([]*models.ModelJobNodes, 0), make([]*models.ModelEdges, 0)

	for _, node := range jobNodes {
		for _, edge := range node.Edges {
			edges = append(edges, &models.ModelEdges{
				Source: node.UUID,
				Target: edge,
			})
		}

		jobs = append(jobs, &models.ModelJobNodes{
			UUID:     node.UUID,
			JobTitle: node.JobTitle,
			JobType:  node.JobType,
			Degree:   node.Degree,
			Data:     node.Data,

			RetryCount:    node.RetryCount,
			RetryInterval: node.RetryInterval,
		})
	}
	return j.InsertJobsToDB(jobs, edges)
}

func (j *JobService) UpdateAndRetriggerJob(uuid string, data json.RawMessage, retryCount *int) error {
	jobs, err := j.jobsRepo.GetJobs(&models.JobFilters{Uuids: []string{uuid}})
	if err != nil {
		return constants.ErrorWrap(constants.ErrJobFetch, err)
	}
	if len(jobs) == 0 {
		return constants.ErrJobNotFound
	}

	job := jobs[0]

	if data != nil {
		job.Data = data
	}
	if retryCount != nil {
		job.RetryCount = *retryCount
	}

	job.Status = constants.OpenJobStatus
	job.RunAfter = time.Now()
	job.JobResponse = nil

	return j.jobsRepo.JobsBulkUpsert([]*models.ModelJobNodes{job})
}

func (j *JobService) GetJobByUUID(uuid string) (*models.ModelJobNodes, error) {
	jobs, err := j.jobsRepo.GetJobs(&models.JobFilters{Uuids: []string{uuid}})
	if err != nil {
		return nil, constants.ErrorWrap(constants.ErrJobFetch, err)
	}
	if len(jobs) == 0 {
		return nil, constants.ErrJobNotFound
	}
	return jobs[0], nil
}
