package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PgContact struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:contacts,alias:c"`

	ID   uint64 `bun:"id,pk,autoincrement" json:"id"`
	UUID string `bun:"uuid,unique" json:"uuid"`

	FirstName   string   `bun:"first_name" json:"first_name"`
	LastName    string   `bun:"last_name" json:"last_name"`
	CompanyID   string   `bun:"company_id" json:"company_id"`
	Email       string   `bun:"email" json:"email"`
	Title       string   `bun:"title" json:"title"`
	Departments []string `bun:"departments,array" json:"departments"`

	MobilePhone string `bun:"mobile_phone" json:"mobile_phone"`
	EmailStatus string `bun:"email_status" json:"email_status"`
	Seniority   string `bun:"seniority" json:"seniority"`
	City        string `bun:"city" json:"city"`
	State       string `bun:"state" json:"state"`
	Country     string `bun:"country" json:"country"`
	LinkedinURL string `bun:"linkedin_url" json:"linkedin_url"`

	FacebookURL     string `bun:"facebook_url" json:"facebook_url"`
	TwitterURL      string `bun:"twitter_url" json:"twitter_url"`
	Website         string `bun:"website" json:"website"`
	WorkDirectPhone string `bun:"work_direct_phone" json:"work_direct_phone"`
	HomePhone       string `bun:"home_phone" json:"home_phone"`
	OtherPhone      string `bun:"other_phone" json:"other_phone"`
	Stage           string `bun:"stage" json:"stage"`

	CreatedAt time.Time  `bun:"created_at,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,nullzero" json:"deleted_at"`
}

func (c *PgContact) SetDB(db *bun.DB) *PgContact {
	c.db = db
	return c
}
