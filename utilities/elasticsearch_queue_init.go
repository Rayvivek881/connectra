package utilities

import (
	"bytes"
	"io"
	"vivek-ray/constants"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

// InitializeElasticsearchQueues initializes the async indexing queues with processors
func InitializeElasticsearchQueues(
	elasticClient *elasticsearch.Client,
) {
	// Company processor
	companyProcessor := func(op *IndexOperation) error {
		switch op.Type {
		case "create", "update":
			// Serialize document
			docBytes, err := SerializeDocument(op.Document)
			if err != nil {
				return err
			}

			// Index document
			response, err := elasticClient.Index(
				constants.CompanyIndex,
				bytes.NewReader(docBytes),
				elasticClient.Index.WithDocumentID(op.DocumentID),
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

		case "delete":
			response, err := elasticClient.Delete(
				constants.CompanyIndex,
				op.DocumentID,
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

		default:
			log.Warn().Str("type", op.Type).Msg("Unknown operation type")
			return nil
		}
	}

	// Contact processor
	contactProcessor := func(op *IndexOperation) error {
		switch op.Type {
		case "create", "update":
			// Serialize document
			docBytes, err := SerializeDocument(op.Document)
			if err != nil {
				return err
			}

			// Index document
			response, err := elasticClient.Index(
				constants.ContactIndex,
				bytes.NewReader(docBytes),
				elasticClient.Index.WithDocumentID(op.DocumentID),
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

		case "delete":
			response, err := elasticClient.Delete(
				constants.ContactIndex,
				op.DocumentID,
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

		default:
			log.Warn().Str("type", op.Type).Msg("Unknown operation type")
			return nil
		}
	}

	// Initialize queues
	InitializeQueues(companyProcessor, contactProcessor)
}
