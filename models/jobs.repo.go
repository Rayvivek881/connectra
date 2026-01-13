package models

import (
	"context"
	"vivek-ray/constants"
	"vivek-ray/utilities"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type JobsStruct struct {
	PgDbClient *bun.DB
}

func JobsRepository(db *bun.DB) JobsSvcRepo {
	return &JobsStruct{
		PgDbClient: db,
	}
}

type JobsFilters struct {
	Uuids    []string
	JobType  string
	Status   []string
	Limit    int
	Retrying bool `default:"false"`
}

func (f *JobsFilters) ToWhereQuery(query *bun.SelectQuery) *bun.SelectQuery {
	if f.JobType != "" {
		query.Where("job_type = ?", f.JobType)
	}
	if len(f.Status) > 0 {
		query.Where("status IN (?)", bun.In(f.Status))
	}
	if f.Retrying {
		query.Where("retry_count > 0")
	}
	if len(f.Uuids) > 0 {
		query.Where("uuid IN (?)", bun.In(f.Uuids))
	}

	limit := utilities.InlineIf(f.Limit > 0, f.Limit, constants.DefaultPageSize).(int)
	return query.Limit(limit)
}

type JobsSvcRepo interface {
	Create(job *ModelJobs) (string, error)
	BulkUpsert(jobs []*ModelJobs) error
	ListByFilters(filters JobsFilters) ([]*ModelJobs, error)
}

func (t *JobsStruct) Create(job *ModelJobs) (string, error) {
	job.UUID = uuid.New().String()

	_, err := t.PgDbClient.NewInsert().
		Model(&job).
		Exec(context.Background())
	return job.UUID, err
}

func (t *JobsStruct) BulkUpsert(jobs []*ModelJobs) error {
	// Deduplicate by UUID - keep the last occurrence
	uniqueJobs := make(map[string]*ModelJobs)
	for _, job := range jobs {
		uniqueJobs[job.UUID] = job
	}

	deduped := make([]*ModelJobs, 0, len(uniqueJobs))
	for _, job := range uniqueJobs {
		deduped = append(deduped, job)
	}

	if len(deduped) == 0 {
		return nil
	}

	_, err := t.PgDbClient.NewInsert().
		Model(&deduped).
		On("CONFLICT(uuid) DO UPDATE").
		Set("status = EXCLUDED.status").
		Set("job_response = EXCLUDED.job_response").
		Set("retry_count = EXCLUDED.retry_count").
		Set("run_after = EXCLUDED.run_after").
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(context.Background())
	return err
}

func (t *JobsStruct) ListByFilters(filters JobsFilters) ([]*ModelJobs, error) {
	var jobs []*ModelJobs
	query := t.PgDbClient.NewSelect().Model(&jobs)
	filters.ToWhereQuery(query)

	err := query.Scan(context.Background())
	return jobs, err
}
