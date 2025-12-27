package models

import (
	"fmt"
	"strings"
	"time"
	"vivek-ray/utilities"

	"github.com/uptrace/bun"
)

type PgContact struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:contacts,alias:c"`

	ID   uint64 `bun:"id,pk,autoincrement" json:"id,omitempty"`
	UUID string `bun:"uuid,unique" json:"uuid,omitempty"`

	FirstName   string   `bun:"first_name" json:"first_name,omitempty"`
	LastName    string   `bun:"last_name" json:"last_name,omitempty"`
	CompanyID   string   `bun:"company_id" json:"company_id,omitempty"`
	Email       string   `bun:"email" json:"email,omitempty"`
	Title       string   `bun:"title" json:"title,omitempty"`
	Departments []string `bun:"departments,array" json:"departments,omitempty"`

	MobilePhone string `bun:"mobile_phone" json:"mobile_phone,omitempty"`
	EmailStatus string `bun:"email_status" json:"email_status,omitempty"`
	Seniority   string `bun:"seniority" json:"seniority,omitempty"`
	City        string `bun:"city" json:"city,omitempty"`
	State       string `bun:"state" json:"state,omitempty"`
	Country     string `bun:"country" json:"country,omitempty"`
	LinkedinURL string `bun:"linkedin_url" json:"linkedin_url,omitempty"`

	FacebookURL     string `bun:"facebook_url" json:"facebook_url,omitempty"`
	TwitterURL      string `bun:"twitter_url" json:"twitter_url,omitempty"`
	Website         string `bun:"website" json:"website,omitempty"`
	WorkDirectPhone string `bun:"work_direct_phone" json:"work_direct_phone,omitempty"`
	HomePhone       string `bun:"home_phone" json:"home_phone,omitempty"`
	OtherPhone      string `bun:"other_phone" json:"other_phone,omitempty"`
	Stage           string `bun:"stage" json:"stage,omitempty"`

	CreatedAt *time.Time `bun:"created_at,nullzero,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:"deleted_at,nullzero" json:"deleted_at,omitempty"`
}

func PgContactFromRowData(row map[string]string, company *PgCompany) *PgContact {
	FirstName, LastName := strings.ToLower(row["first_name"]), strings.ToLower(row["last_name"])
	LinkedinURL := strings.ToLower(row["person_linkedin_url"])
	server_time := time.Now()

	ContactUUID := utilities.GenerateUUID5(fmt.Sprintf("%s%s%s", FirstName, LastName, LinkedinURL))
	return &PgContact{
		UUID: ContactUUID,

		FirstName:   FirstName,
		LastName:    LastName,
		CompanyID:   company.UUID,
		Email:       strings.ToLower(row["email"]),
		Title:       row["title"],
		Departments: utilities.SplitAndTrim(row["departments"], ","),

		MobilePhone: utilities.GetCleanedPhoneNumber(row["mobile_phone"]),
		EmailStatus: strings.ToLower(row["email_status"]),
		Seniority:   strings.ToLower(row["seniority"]),
		City:        strings.ToLower(row["city"]),
		State:       strings.ToLower(row["state"]),
		Country:     strings.ToLower(row["country"]),
		LinkedinURL: LinkedinURL,

		FacebookURL:     row["facebook_url"],
		TwitterURL:      row["twitter_url"],
		WorkDirectPhone: utilities.GetCleanedPhoneNumber(row["work_direct_phone"]),
		HomePhone:       utilities.GetCleanedPhoneNumber(row["home_phone"]),
		OtherPhone:      utilities.GetCleanedPhoneNumber(row["other_phone"]),
		Stage:           strings.ToLower(row["stage"]),
		CreatedAt:       &server_time,
		UpdatedAt:       &server_time,
	}
}

func (c *PgContact) SetDB(db *bun.DB) *PgContact {
	c.db = db
	return c
}
