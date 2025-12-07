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
