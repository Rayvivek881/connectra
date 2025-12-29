package models

import (
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticCompany struct {
	client *elasticsearch.Client

	Id string `json:"id"`

	Name             string   `json:"name"`              // text search
	EmployeesCount   int64    `json:"employees_count"`   // number search
	Industries       []string `json:"industries"`        // keyword search
	Keywords         []string `json:"keywords"`          // keyword search
	Address          string   `json:"address"`           // text search
	AnnualRevenue    int64    `json:"annual_revenue"`    // number search
	TotalFunding     int64    `json:"total_funding"`     // number search
	Technologies     []string `json:"technologies"`      // keyword search
	City             string   `json:"city"`              // text search
	State            string   `json:"state"`             // text search
	Country          string   `json:"country"`           // text search
	LinkedinURL      string   `json:"linkedin_url"`      // text search
	Website          string   `json:"website"`           // text search
	NormalizedDomain string   `json:"normalized_domain"` // text search

	CreatedAt time.Time `json:"created_at"` // date search
}

type ElasticCompanySearchResponse struct {
	Hits struct {
		Hits []struct {
			Source ElasticCompany `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (c *ElasticCompany) SetClient(client *elasticsearch.Client) *ElasticCompany {
	c.client = client
	return c
}
