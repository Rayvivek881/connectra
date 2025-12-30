package jobs

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	commonService "vivek-ray/modules/common/service"
	companyService "vivek-ray/modules/companies/service"
	contactService "vivek-ray/modules/contacts/service"
	"vivek-ray/utilities"

	"github.com/rs/zerolog/log"
)

func InsertCsvToDb(fileStream *io.ReadCloser) error {
	csvReader, batchUpsertService := csv.NewReader(*fileStream), commonService.NewBatchUpsertService()
	headers, err := csvReader.Read()
	if err != nil {
		return err
	}
	batchSize := conf.JobConfig.BatchSize
	batch := make([]map[string]string, 0, batchSize)

	for {
		row, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		batch = append(batch, utilities.CsvRowToMap(headers, row))
		if len(batch) >= batchSize {
			if err := batchUpsertService.ProcessBatchUpsert(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		return batchUpsertService.ProcessBatchUpsert(batch)
	}
	return nil
}

func ProcessInsertCsvFile(job *models.ModelJobs) error {
	var jobData utilities.InsertFileJobData
	if err := json.Unmarshal(job.Data, &jobData); err != nil {
		return err
	}
	if jobData.FileS3Bucket == "" {
		jobData.FileS3Bucket = conf.S3StorageConfig.S3Bucket
	}
	fileStream, err := connections.S3Connection.ReadFileStream(
		context.Background(),
		jobData.FileS3Bucket,
		jobData.FileS3Key,
	)
	if err != nil {
		return err
	}
	defer fileStream.Close()
	return InsertCsvToDb(&fileStream)
}

func ExportContactsCsvToStream(writer *io.PipeWriter, vql utilities.VQLQuery) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	vql.CompanyConfig = nil
	service := contactService.NewContactService([]*models.ModelFilter{})

	if err := csvWriter.Write(vql.SelectColumns); err != nil {
		return err
	}
	for {
		contacts, err := service.ListByFilters(vql)
		if err != nil {
			return err
		}
		if len(contacts) == 0 {
			break
		}

		for _, contact := range contacts {
			row := utilities.StructToCsvSlice(contact.PgContact, vql.SelectColumns)
			if err := csvWriter.Write(row); err != nil {
				return err
			}
		}

		vql.SearchAfter = contacts[len(contacts)-1].SearchAfter
		csvWriter.Flush()
	}
	return nil
}

func ExportCompaniesCsvToStream(writer *io.PipeWriter, vql utilities.VQLQuery) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()
	if err := csvWriter.Write(vql.SelectColumns); err != nil {
		return err
	}

	service := companyService.NewCompanyService([]*models.ModelFilter{})
	for {
		companies, err := service.ListByFilters(vql)
		if err != nil {
			return err
		}
		if len(companies) == 0 {
			break
		}
		for _, company := range companies {
			row := utilities.StructToCsvSlice(company, vql.SelectColumns)
			if err := csvWriter.Write(row); err != nil {
				return err
			}
		}
		vql.SearchAfter = companies[len(companies)-1].SearchAfter
		csvWriter.Flush()
	}
	return nil
}

func ExportCsvToStream(writer *io.PipeWriter, jobData utilities.ExportFileJobData) error {
	vql := jobData.VQL
	vql.OrderBy = []utilities.FilterOrder{{OrderBy: "created_at", OrderDirection: "desc"}}
	vql.Limit = conf.JobConfig.BatchSize

	if len(vql.SelectColumns) == 0 {
		return errors.New("select columns are required")
	}

	switch jobData.Service {
	case constants.ContactsService:
		return ExportContactsCsvToStream(writer, vql)
	case constants.CompaniesService:
		return ExportCompaniesCsvToStream(writer, vql)
	default:
		return errors.New("invalid service")
	}
}

func ProcessExportCsvFile(job *models.ModelJobs) error {
	var jobData utilities.ExportFileJobData
	if err := json.Unmarshal(job.Data, &jobData); err != nil {
		return err
	}
	if jobData.FileS3Bucket == "" {
		jobData.FileS3Bucket = conf.S3StorageConfig.S3Bucket
	}
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		if err := ExportCsvToStream(writer, jobData); err != nil {
			log.Error().Err(err).Msg("Failed to export csv to stream")
		}
	}()

	s3Key := fmt.Sprintf("%s/%s.csv", conf.S3StorageConfig.S3UploadFilePath, job.UUID)

	if err := connections.S3Connection.WriteFileStream(context.Background(), jobData.FileS3Bucket, s3Key, reader); err != nil {
		return err
	}

	job.AddS3Key(s3Key)
	return nil
}
