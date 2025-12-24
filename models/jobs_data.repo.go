package models

import (
	"context"
	"time"
	"vivek-ray/constants"
	"vivek-ray/utilities"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type JobsDataStruct struct {
	PgDbClient *bun.DB
}

func JobsDataRepository(db *bun.DB) JobsDataSvcRepo {
	return &JobsDataStruct{
		PgDbClient: db,
	}
}

type JobsDataFilters struct {
	Uuids      []string
	JobType    string
	Status     string
	RetryAfter *time.Time
	limit      int
}

func (f *JobsDataFilters) ToWhereQuery(query *bun.SelectQuery) *bun.SelectQuery {
	if len(f.Uuids) > 0 {
		query.Where("uuid IN (?)", bun.In(f.Uuids))
	}
	if f.JobType != "" {
		query.Where("job_type = ?", f.JobType)
	}
	if f.Status != "" {
		query.Where("status = ?", f.Status)
	}

	if f.RetryAfter != nil {
		query.Where("retry_after <= ?", f.RetryAfter)
	}

	limit := utilities.InlineIf(f.limit > 0, f.limit, constants.DefaultPageSize).(int)
	return query.Limit(limit)
}

type JobsDataSvcRepo interface {
	Create(job *ModelJobsData) (string, error)
	Update(uuid string, job *ModelJobsData) error
	ListByFilters(filters JobsDataFilters) ([]*ModelJobsData, error)
}

func (t *JobsDataStruct) Create(job *ModelJobsData) (string, error) {
	job.UUID = uuid.New().String()

	_, err := t.PgDbClient.NewInsert().
		Model(&job).
		Exec(context.Background())
	return job.UUID, err
}

func (t *JobsDataStruct) Update(uuid string, job *ModelJobsData) error {
	_, err := t.PgDbClient.NewUpdate().
		Model(&job).
		Where("uuid = ?", uuid).
		Exec(context.Background())
	return err
}

func (t *JobsDataStruct) ListByFilters(filters JobsDataFilters) ([]*ModelJobsData, error) {
	var jobs []*ModelJobsData
	query := t.PgDbClient.NewSelect().Model(&jobs)
	filters.ToWhereQuery(query)

	err := query.Scan(context.Background())
	return jobs, err
}
