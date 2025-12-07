package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PgCompany struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:companies,alias:cp"`

	ID   uint64 `bun:"id,pk,autoincrement" json:"id"`
	UUID string `bun:"uuid,notnull,unique" json:"uuid"`

	Name             string   `bun:"name" json:"name"`
	EmployeesCount   int64    `bun:"employees_count" json:"employees_count"`
	Industries       []string `bun:"industries,array" json:"industries"`
	Keywords         []string `bun:"keywords,array" json:"keywords"`
	Address          string   `bun:"address" json:"address"`
	AnnualRevenue    int64    `bun:"annual_revenue" json:"annual_revenue"`
	TotalFunding     int64    `bun:"total_funding" json:"total_funding"`
	Technologies     []string `bun:"technologies,array" json:"technologies"`
	City             string   `bun:"city" json:"city"`
	State            string   `bun:"state" json:"state"`
	Country          string   `bun:"country" json:"country"`
	LinkedinURL      string   `bun:"linkedin_url" json:"linkedin_url"`
	Website          string   `bun:"website" json:"website"`
	NormalizedDomain string   `bun:"normalized_domain" json:"normalized_domain"`

	FacebookURL          string `bun:"facebook_url" json:"facebook_url"`
	TwitterURL           string `bun:"twitter_url" json:"twitter_url"`
	CompanyNameForEmails string `bun:"company_name_for_emails" json:"company_name_for_emails"`
	PhoneNumber          string `bun:"phone_number" json:"phone_number"`
	LatestFunding        string `bun:"latest_funding" json:"latest_funding"`
	LatestFundingAmount  int64  `bun:"latest_funding_amount" json:"latest_funding_amount"`
	LastRaisedAt         string `bun:"last_raised_at" json:"last_raised_at"`

	CreatedAt time.Time  `bun:"created_at,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,nullzero" json:"deleted_at"`
}

func (c *PgCompany) SetDB(db *bun.DB) *PgCompany {
	c.db = db
	return c
}
