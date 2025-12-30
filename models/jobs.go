package models

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type JobResponseData struct {
	RuntimeErrors []string `json:"runtime_errors,omitempty"`
	Messages      string   `json:"messages,omitempty"`
	S3Key         string   `json:"s3_key,omitempty"`
}

type ModelJobs struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:jobs,alias:j"`

	Id          uint64          `bun:"id,pk,autoincrement" json:"id"`
	UUID        string          `bun:"uuid,notnull,unique" json:"uuid"`
	JobType     string          `bun:"job_type,notnull" json:"job_type"`
	Data        json.RawMessage `bun:"data,type:jsonb" json:"data"`
	Status      string          `bun:"status,notnull" json:"status"`
	JobResponse json.RawMessage `bun:"job_response,type:jsonb" json:"job_response"`

	RetryCount    int        `bun:"retry_count,notnull,default:0" json:"retry_count"`
	RetryInterval int        `bun:"retry_interval,notnull,default:0" json:"retry_interval"`
	RunAfter      *time.Time `bun:"run_after,nullzero" json:"run_after"`

	CreatedAt *time.Time `bun:"created_at,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at"`
}

func (m *ModelJobs) SetDB(db *bun.DB) *ModelJobs {
	m.db = db
	return m
}

func (m *ModelJobs) AddRuntimeError(errMsg string) {
	var resp JobResponseData
	if len(m.JobResponse) > 0 {
		json.Unmarshal(m.JobResponse, &resp)
	}
	resp.RuntimeErrors = append(resp.RuntimeErrors, errMsg)
	m.JobResponse, _ = json.Marshal(resp)
}

func (m *ModelJobs) AddMessage(message string) {
	var resp JobResponseData
	if len(m.JobResponse) > 0 {
		json.Unmarshal(m.JobResponse, &resp)
	}
	resp.Messages = message
	m.JobResponse, _ = json.Marshal(resp)
}

func (m *ModelJobs) AddS3Key(s3Key string) {
	var resp JobResponseData
	if len(m.JobResponse) > 0 {
		json.Unmarshal(m.JobResponse, &resp)
	}
	resp.S3Key = s3Key
	m.JobResponse, _ = json.Marshal(resp)
}
