package main

import (
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/uptrace/bun"
)

// ============================================
// PostgreSQL Models
// ============================================

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

	// Extra fields for PG only
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

	// Extra fields for PG only
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

type PgFilter struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:filters,alias:cf"`

	Id           uint64     `bun:"id,pk,autoincrement" json:"id"`
	Key          string     `bun:"key" json:"key"`
	Service      string     `bun:"service" json:"service"`
	DisplayName  string     `bun:"display_name,notnull" json:"display_name"`
	DirectDrived bool       `bun:"direct_drived,nullzero" json:"direct_drived"`
	DeletedAt    *time.Time `bun:"deleted_at,nullzero" json:"deleted_at"`
}

func (f *PgFilter) SetDB(db *bun.DB) *PgFilter {
	f.db = db
	return f
}

type PgFilterData struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:filters_data,alias:cfd"`

	Id           uint64     `bun:"id,pk,autoincrement" json:"id"`
	FilterKey    string     `bun:"filter_key,notnull" json:"filter_key"`
	Service      string     `bun:"service,notnull" json:"service"`
	DisplayValue string     `bun:"display_value,notnull" json:"display_value"`
	Value        string     `bun:"value,nullzero" json:"value"`
	DeletedAt    *time.Time `bun:"deleted_at,nullzero" json:"deleted_at"`
}

func (fd *PgFilterData) SetDB(db *bun.DB) *PgFilterData {
	fd.db = db
	return fd
}

// ============================================
// Elasticsearch Models
// ============================================

type ElasticCompany struct {
	client *elasticsearch.Client

	Id string `json:"id"`

	Name             string   `json:"name"`
	EmployeesCount   int64    `json:"employees_count"`
	Industries       []string `json:"industries"`
	Keywords         []string `json:"keywords"`
	Address          string   `json:"address"`
	AnnualRevenue    int64    `json:"annual_revenue"`
	TotalFunding     int64    `json:"total_funding"`
	Technologies     []string `json:"technologies"`
	City             string   `json:"city"`
	State            string   `json:"state"`
	Country          string   `json:"country"`
	LinkedinURL      string   `json:"linkedin_url"`
	Website          string   `json:"website"`
	NormalizedDomain string   `json:"normalized_domain"`

	CreatedAt time.Time `json:"created_at"`
}

func (c *ElasticCompany) SetClient(client *elasticsearch.Client) *ElasticCompany {
	c.client = client
	return c
}

type ElasticContact struct {
	client *elasticsearch.Client

	Id string `json:"id"`

	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	CompanyID   string   `json:"company_id"`
	Email       string   `json:"email"`
	Title       string   `json:"title"`
	Departments []string `json:"departments"`

	MobilePhone string `json:"mobile_phone"`
	EmailStatus string `json:"email_status"`
	Seniority   string `json:"seniority"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	LinkedinURL string `json:"linkedin_url"`

	// company details
	CompanyName             string   `json:"company_name"`              // text search
	CompanyEmployeesCount   int64    `json:"company_employees_count"`   // number search
	CompanyIndustries       []string `json:"company_industries"`        // keyword search
	CompanyKeywords         []string `json:"company_keywords"`          // keyword search
	CompanyAddress          string   `json:"company_address"`           // text search
	CompanyAnnualRevenue    int64    `json:"company_annual_revenue"`    // number search
	CompanyTotalFunding     int64    `json:"company_total_funding"`     // number search
	CompanyTechnologies     []string `json:"company_technologies"`      // keyword search
	CompanyCity             string   `json:"company_city"`              // text search
	CompanyState            string   `json:"company_state"`             // text search
	CompanyCountry          string   `json:"company_country"`           // text search
	CompanyLinkedinURL      string   `json:"company_linkedin_url"`      // text search
	CompanyWebsite          string   `json:"company_website"`           // text search
	CompanyNormalizedDomain string   `json:"company_normalized_domain"` // text search

	CreatedAt time.Time `json:"created_at"`
}

func (c *ElasticContact) SetClient(client *elasticsearch.Client) *ElasticContact {
	c.client = client
	return c
}

// ============================================
// Batch Types for Producer-Consumer
// ============================================

// CompanyBatch holds a batch of companies and their contacts
type CompanyBatch struct {
	BatchNum    int
	PgCompanies []PgCompany
	PgContacts  []PgContact
	EsCompanies []ElasticCompany
	EsContacts  []ElasticContact
}