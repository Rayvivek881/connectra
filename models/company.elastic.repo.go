package models

import (
	"bytes"
	"encoding/json"
	"io"
	"vivek-ray/constants"
	"vivek-ray/utilities"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

type ElasticCompanyStruct struct {
	ElasticClient *elasticsearch.Client
}

func ElasticCompanyRepository(client *elasticsearch.Client) ElasticCompanySvcRepo {
	return &ElasticCompanyStruct{
		ElasticClient: client,
	}
}

type ElasticCompanySvcRepo interface {
	ListByQueryMap(query map[string]any) ([]*ElasticCompanySearchHit, error)
	CountByQueryMap(query map[string]any) (int64, error)
	BulkUpsert(companies []*ElasticCompany) (int64, error)
	Create(company *ElasticCompany) error
	Update(company *ElasticCompany) error
	Delete(uuid string) error
}

func (t *ElasticCompanyStruct) ListByQueryMap(query map[string]any) ([]*ElasticCompanySearchHit, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.ElasticClient.Search(
		t.ElasticClient.Search.WithIndex(constants.CompanyIndex),
		t.ElasticClient.Search.WithBody(queryReader),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, constants.ElasticsearchError(response.StatusCode, string(bodyBytes))
	}

	var searchResponse ElasticCompanySearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	return searchResponse.Hits.Hits, nil
}

func (t *ElasticCompanyStruct) CountByQueryMap(query map[string]any) (int64, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.ElasticClient.Count(
		t.ElasticClient.Count.WithIndex(constants.CompanyIndex),
		t.ElasticClient.Count.WithBody(queryReader),
	)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, constants.ElasticsearchError(response.StatusCode, string(bodyBytes))
	}

	var countResponse utilities.ElasticCount
	if err := json.NewDecoder(response.Body).Decode(&countResponse); err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

func (t *ElasticCompanyStruct) BulkUpsert(companies []*ElasticCompany) (int64, error) {
	var buf bytes.Buffer
	for _, company := range companies {
		meta := map[string]any{
			"index": map[string]any{
				"_index": constants.CompanyIndex,
				"_id":    company.UUID,
			},
		}
		if utilities.AddToBuffer(&buf, meta) != nil || utilities.AddToBuffer(&buf, company) != nil {
			log.Error().Msgf("Failed to add company to buffer: %v", company.UUID)
			continue
		}
	}

	response, err := t.ElasticClient.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, constants.ElasticsearchBulkError(response.StatusCode, string(bodyBytes))
	}

	return int64(len(companies)), nil
}

func (t *ElasticCompanyStruct) Create(company *ElasticCompany) error {
	companyJson, err := json.Marshal(company)
	if err != nil {
		return err
	}

	response, err := t.ElasticClient.Index(
		constants.CompanyIndex,
		bytes.NewReader(companyJson),
		t.ElasticClient.Index.WithDocumentID(company.UUID),
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return constants.ElasticsearchError(response.StatusCode, string(bodyBytes))
	}

	return nil
}

func (t *ElasticCompanyStruct) Update(company *ElasticCompany) error {
	// In Elasticsearch, update is the same as index (upsert behavior)
	return t.Create(company)
}

func (t *ElasticCompanyStruct) Delete(uuid string) error {
	response, err := t.ElasticClient.Delete(
		constants.CompanyIndex,
		uuid,
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return constants.ElasticsearchError(response.StatusCode, string(bodyBytes))
	}

	return nil
}
