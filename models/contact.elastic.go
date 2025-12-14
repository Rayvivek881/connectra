package models

import (
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticContact struct {
	client *elasticsearch.Client

	Id string `json:"id"`

	FirstName   string   `json:"first_name"`  // text search
	LastName    string   `json:"last_name"`   // text search
	CompanyID   string   `json:"company_id"`  // keyword search
	Email       string   `json:"email"`       // keyword search
	Title       string   `json:"title"`       // text search
	Departments []string `json:"departments"` // keyword search

	MobilePhone string `json:"mobile_phone"` // keyword search
	EmailStatus string `json:"email_status"` // keyword search
	Seniority   string `json:"seniority"`    // keyword search
	City        string `json:"city"`         // text search
	State       string `json:"state"`        // text search
	Country     string `json:"country"`      // text search
	LinkedinURL string `json:"linkedin_url"` // text search

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

type ElasticContactSearchResponse struct {
	Hits struct {
		Hits []struct {
			Source ElasticContact `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (c *ElasticContact) SetClient(client *elasticsearch.Client) *ElasticContact {
	c.client = client
	return c
}
