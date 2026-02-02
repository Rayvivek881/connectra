package models

import (
	"time"

	"github.com/opensearch-project/opensearch-go/v2"
)

type OpenSearchCompany struct {
	client *opensearch.Client

	UUID string `json:"uuid"`

	Name             string   `json:"name"`              // text search
	EmployeesCount   int64    `json:"employees_count"`   // number search
	Industries       []string `json:"industries"`        // keyword search
	Keywords         []string `json:"keywords"`          // keyword search
	Address          string   `json:"address"`           // text search
	AnnualRevenue    int64    `json:"annual_revenue"`    // number search
	TotalFunding     int64    `json:"total_funding"`     // number search
	Technologies     []string `json:"technologies"`       // keyword search
	City             string   `json:"city"`              // text search
	State            string   `json:"state"`              // text search
	Country          string   `json:"country"`           // text search
	LinkedinURL      string   `json:"linkedin_url"`      // text search
	Website          string   `json:"website"`            // text search
	NormalizedDomain string   `json:"normalized_domain"` // text search

	CreatedAt time.Time `json:"created_at"` // date search
}

func OpenSearchCompanyFromRawData(company *PgCompany) *OpenSearchCompany {
	serverTime := time.Now()
	return &OpenSearchCompany{
		UUID:             company.UUID,
		Name:             company.Name,
		EmployeesCount:   company.EmployeesCount,
		Industries:       company.Industries,
		Keywords:         company.Keywords,
		Address:          company.Address,
		AnnualRevenue:    company.AnnualRevenue,
		TotalFunding:     company.TotalFunding,
		Technologies:     company.Technologies,
		City:             company.City,
		State:            company.State,
		Country:          company.Country,
		LinkedinURL:      company.LinkedinURL,
		Website:          company.Website,
		NormalizedDomain: company.NormalizedDomain,
		CreatedAt:        serverTime,
	}
}

type OpenSearchCompanySearchHit struct {
	Company OpenSearchCompany `json:"_source"`
	Cursor  []string          `json:"sort,omitempty"`
}
type OpenSearchCompanySearchResponse struct {
	Hits struct {
		Hits []*OpenSearchCompanySearchHit `json:"hits"`
	} `json:"hits"`
}

func (c *OpenSearchCompany) SetClient(client *opensearch.Client) *OpenSearchCompany {
	c.client = client
	return c
}
