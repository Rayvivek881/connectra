package models

import (
	"fmt"
	"strings"
	"time"
	"vivek-ray/utilities"

	"github.com/uptrace/bun"
)

type PgCompany struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:companies,alias:cp"`

	ID   uint64 `bun:"id,pk,autoincrement" json:"id,omitempty"`
	UUID string `bun:"uuid,notnull,unique" json:"uuid,omitempty"`

	Name             string   `bun:"name" json:"name,omitempty"`
	EmployeesCount   int64    `bun:"employees_count" json:"employees_count,omitempty"`
	Industries       []string `bun:"industries,array" json:"industries,omitempty"`
	Keywords         []string `bun:"keywords,array" json:"keywords,omitempty"`
	Address          string   `bun:"address" json:"address,omitempty"`
	AnnualRevenue    int64    `bun:"annual_revenue" json:"annual_revenue,omitempty"`
	TotalFunding     int64    `bun:"total_funding" json:"total_funding,omitempty"`
	Technologies     []string `bun:"technologies,array" json:"technologies,omitempty"`
	City             string   `bun:"city" json:"city,omitempty"`
	State            string   `bun:"state" json:"state,omitempty"`
	Country          string   `bun:"country" json:"country,omitempty"`
	LinkedinURL      string   `bun:"linkedin_url" json:"linkedin_url,omitempty"`
	Website          string   `bun:"website" json:"website,omitempty"`
	NormalizedDomain string   `bun:"normalized_domain" json:"normalized_domain,omitempty"`

	FacebookURL          string `bun:"facebook_url" json:"facebook_url,omitempty"`
	TwitterURL           string `bun:"twitter_url" json:"twitter_url,omitempty"`
	CompanyNameForEmails string `bun:"company_name_for_emails" json:"company_name_for_emails,omitempty"`
	PhoneNumber          string `bun:"phone_number" json:"phone_number,omitempty"`
	LatestFunding        string `bun:"latest_funding" json:"latest_funding,omitempty"`
	LatestFundingAmount  int64  `bun:"latest_funding_amount" json:"latest_funding_amount,omitempty"`
	LastRaisedAt         string `bun:"last_raised_at" json:"last_raised_at,omitempty"`

	CreatedAt *time.Time `bun:"created_at,nullzero,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:"deleted_at,nullzero" json:"deleted_at,omitempty"`
}

func PgCompanyFromRawData(row map[string]string) *PgCompany {
	companyName, email := row["company"], strings.ToLower(row["email"])
	server_time, linkedinURL := time.Now(), strings.ToLower(row["company_linkedin_url"])
	var normalizedDomain string
	if _, domain, found := strings.Cut(email, "@"); found {
		normalizedDomain = domain
	}

	CompanyUUID := utilities.GenerateUUID5(fmt.Sprintf("%s%s", strings.ToLower(companyName), linkedinURL))
	return &PgCompany{
		UUID: CompanyUUID,

		Name:             companyName,
		EmployeesCount:   utilities.StringToInt64(row["employees"]),
		Industries:       utilities.SplitAndTrim(row["industry"], ","),
		Keywords:         utilities.SplitAndTrim(row["keywords"], ","),
		Address:          row["company_address"],
		AnnualRevenue:    utilities.StringToInt64(row["annual_revenue"]),
		TotalFunding:     utilities.StringToInt64(row["total_funding"]),
		Technologies:     utilities.SplitAndTrim(row["technologies"], ","),
		Website:          row["website"],
		LinkedinURL:      linkedinURL,
		City:             strings.ToLower(row["company_city"]),
		State:            strings.ToLower(row["company_state"]),
		Country:          strings.ToLower(row["company_country"]),
		NormalizedDomain: normalizedDomain,

		FacebookURL:          row["facebook_url"],
		TwitterURL:           row["twitter_url"],
		CompanyNameForEmails: row["company_name_for_emails"],
		PhoneNumber:          row["company_phone"],
		LatestFunding:        row["latest_funding"],
		LatestFundingAmount:  utilities.StringToInt64(row["latest_funding_amount"]),
		LastRaisedAt:         row["last_raised_at"],

		CreatedAt: &server_time,
		UpdatedAt: &server_time,
	}
}

func (c *PgCompany) SetDB(db *bun.DB) *PgCompany {
	c.db = db
	return c
}
