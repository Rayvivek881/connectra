package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// ============================================
// Configuration
// ============================================

const (
	// Total records
	TotalCompanies     = 1000000 // 1 million companies
	ContactsPerCompany = 5       // 5 contacts per company
	TotalContacts      = TotalCompanies * ContactsPerCompany

	// Batch sizes
	CompanyBatchSize = 1000 // Companies per batch
	ContactBatchSize = 1000 // Contacts per batch (will be processed as 5000 per company batch)

	// Concurrency
	NumProducers      = 2       // Number of producer goroutines
	NumPgConsumers    = 4       // Number of PostgreSQL consumer goroutines
	NumEsConsumers    = 6       // Number of Elasticsearch consumer goroutines
	ChannelBufferSize = 1000000 // Buffer size for channels

	// Database settings
	PgDSN = "postgres://postgres:bc1q9p333c20qemhvjmadylvsmht5xj2z5fs27pu3u@98.81.200.121:5432/production"
	// PgDSN = "postgres://samosa:12345@localhost:5432/postgres"

	// Elasticsearch settings
	// EsURL      = "http://localhost:9200"
	EsURL      = "http://54.198.202.181:9200"
	EsUsername = "elastic"
	EsPassword = "afbsdcvgdbkjgyhbnjs127gh"

	// Index names
	EsCompanyIndex = "companies_index"
	EsContactIndex = "contacts_index"
)

// ============================================
// Statistics Tracking
// ============================================

type Stats struct {
	PgCompaniesIndexed int64
	PgCompaniesFailed  int64
	PgContactsIndexed  int64
	PgContactsFailed   int64
	EsCompaniesIndexed int64
	EsCompaniesFailed  int64
	EsContactsIndexed  int64
	EsContactsFailed   int64
	BatchesCompleted   int64
	TotalBatches       int64
}

func (s *Stats) AddPgCompanies(success, failed int64) {
	atomic.AddInt64(&s.PgCompaniesIndexed, success)
	atomic.AddInt64(&s.PgCompaniesFailed, failed)
}

func (s *Stats) AddPgContacts(success, failed int64) {
	atomic.AddInt64(&s.PgContactsIndexed, success)
	atomic.AddInt64(&s.PgContactsFailed, failed)
}

func (s *Stats) AddEsCompanies(success, failed int64) {
	atomic.AddInt64(&s.EsCompaniesIndexed, success)
	atomic.AddInt64(&s.EsCompaniesFailed, failed)
}

func (s *Stats) AddEsContacts(success, failed int64) {
	atomic.AddInt64(&s.EsContactsIndexed, success)
	atomic.AddInt64(&s.EsContactsFailed, failed)
}

func (s *Stats) IncrementBatches() {
	atomic.AddInt64(&s.BatchesCompleted, 1)
}

func (s *Stats) Print() {
	fmt.Printf("\n========== Current Stats ==========\n")
	fmt.Printf("Batches: %d/%d\n", atomic.LoadInt64(&s.BatchesCompleted), s.TotalBatches)
	fmt.Printf("PG Companies - Success: %d, Failed: %d\n",
		atomic.LoadInt64(&s.PgCompaniesIndexed), atomic.LoadInt64(&s.PgCompaniesFailed))
	fmt.Printf("PG Contacts  - Success: %d, Failed: %d\n",
		atomic.LoadInt64(&s.PgContactsIndexed), atomic.LoadInt64(&s.PgContactsFailed))
	fmt.Printf("ES Companies - Success: %d, Failed: %d\n",
		atomic.LoadInt64(&s.EsCompaniesIndexed), atomic.LoadInt64(&s.EsCompaniesFailed))
	fmt.Printf("ES Contacts  - Success: %d, Failed: %d\n",
		atomic.LoadInt64(&s.EsContactsIndexed), atomic.LoadInt64(&s.EsContactsFailed))
	fmt.Printf("====================================\n\n")
}

// ============================================
// Helper Functions
// ============================================

func toLowerCase(s string) string {
	return strings.ToLower(s)
}

func weightedChoice(weights map[string]float64) string {
	total := 0.0
	for _, w := range weights {
		total += w
	}

	r := rand.Float64() * total
	cumulative := 0.0
	for k, w := range weights {
		cumulative += w
		if r <= cumulative {
			return k
		}
	}

	// Return first key as fallback
	for k := range weights {
		return k
	}
	return ""
}

func randomSample(slice []string, n int) []string {
	if n > len(slice) {
		n = len(slice)
	}
	indices := rand.Perm(len(slice))[:n]
	result := make([]string, n)
	for i, idx := range indices {
		result[i] = slice[idx]
	}
	return result
}

func randomInt64InRange(min, max int64) int64 {
	if max <= min {
		return min
	}
	return min + rand.Int63n(max-min+1)
}

// ============================================
// Company Generation
// ============================================

func generateCompanyName() string {
	nameType := rand.Intn(4)
	switch nameType {
	case 0: // prefix_suffix
		return CompanyNamePrefixes[rand.Intn(len(CompanyNamePrefixes))] +
			CompanyNameSuffixes[rand.Intn(len(CompanyNameSuffixes))]
	case 1: // word_suffix
		return CompanyNameWords[rand.Intn(len(CompanyNameWords))] + " " +
			CompanyNameSuffixes[rand.Intn(len(CompanyNameSuffixes))]
	case 2: // two_words
		return CompanyNameWords[rand.Intn(len(CompanyNameWords))] + " " +
			CompanyNameWords[rand.Intn(len(CompanyNameWords))]
	default: // single_word
		return CompanyNameWords[rand.Intn(len(CompanyNameWords))]
	}
}

func generateDomainFromName(companyName string) string {
	domainName := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(companyName, " ", ""), "-", ""))
	if rand.Float64() < 0.2 {
		domainName += fmt.Sprintf("%d", rand.Intn(999)+1)
	}
	extension := DomainExtensions[rand.Intn(len(DomainExtensions))]
	return domainName + extension
}

func generateCompanyRecord() (PgCompany, ElasticCompany) {
	country := Countries[rand.Intn(len(Countries))]

	var state string
	if states, ok := StatesByCountry[country]; ok && len(states) > 0 {
		state = states[rand.Intn(len(states))]
	}

	var city string
	if cities, ok := CitiesByCountry[country]; ok && len(cities) > 0 {
		city = cities[rand.Intn(len(cities))]
	} else {
		city = gofakeit.City()
	}

	// Generate company name and related fields
	companyName := generateCompanyName()
	domain := generateDomainFromName(companyName)
	normalizedDomain := strings.Split(domain, ".")[0]
	website := "https://www." + domain
	linkedinURL := fmt.Sprintf("https://linkedin.com/company/%s-%d", normalizedDomain, rand.Intn(900000)+100000)

	// Determine company size and funding stage
	companySize := weightedChoice(CompanySizeWeights)
	fundingStage := weightedChoice(FundingStageWeights)

	// Generate employees count based on company size
	empRange := EmployeeCountRanges[companySize]
	employeesCount := randomInt64InRange(empRange.Min, empRange.Max)

	// Generate annual revenue based on company size
	revRange := AnnualRevenueRanges[companySize]
	annualRevenue := randomInt64InRange(revRange.Min, revRange.Max)

	// Generate total funding based on funding stage
	fundRange := TotalFundingRanges[fundingStage]
	var totalFunding int64
	if fundRange.Max > 0 {
		totalFunding = randomInt64InRange(fundRange.Min, fundRange.Max)
	}

	// Generate arrays
	numIndustries := rand.Intn(3) + 1
	industries := randomSample(Industries, numIndustries)

	numKeywords := rand.Intn(4) + 2
	keywords := randomSample(Keywords, numKeywords)

	numTechnologies := rand.Intn(8) + 3
	technologies := randomSample(Technologies, numTechnologies)

	// Generate address
	address := fmt.Sprintf("%s, %s, %s %s, %s",
		gofakeit.StreetName()+" "+gofakeit.StreetSuffix(),
		city, state, gofakeit.Zip(), country)

	// Generate created_at date (within last 5 years)
	daysAgo := rand.Intn(1825)
	createdAt := time.Now().AddDate(0, 0, -daysAgo)

	// Generate deterministic UUID based on name, NormalizedDomain, LinkedinURL
	uuidSeed := fmt.Sprintf("%s%s%s", companyName, normalizedDomain, linkedinURL)
	companyUUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(uuidSeed)).String()

	// Generate extra PG fields
	facebookURL := fmt.Sprintf("https://facebook.com/%s", normalizedDomain)
	twitterURL := fmt.Sprintf("https://twitter.com/%s", normalizedDomain)
	companyNameForEmails := strings.ToLower(strings.ReplaceAll(companyName, " ", ""))
	phoneNumber := gofakeit.Phone()
	latestFunding := FundingStages[rand.Intn(len(FundingStages))]
	latestFundingAmount := randomInt64InRange(100000, 50000000)
	lastRaisedAt := time.Now().AddDate(0, 0, -rand.Intn(365)).Format("2006-01-02")

	pgCompany := PgCompany{
		UUID:                 companyUUID,
		Name:                 companyName,
		EmployeesCount:       employeesCount,
		Industries:           industries,
		Keywords:             keywords,
		Address:              address,
		AnnualRevenue:        annualRevenue,
		TotalFunding:         totalFunding,
		Technologies:         technologies,
		City:                 city,
		State:                state,
		Country:              country,
		LinkedinURL:          linkedinURL,
		Website:              website,
		NormalizedDomain:     normalizedDomain,
		FacebookURL:          facebookURL,
		TwitterURL:           twitterURL,
		CompanyNameForEmails: companyNameForEmails,
		PhoneNumber:          phoneNumber,
		LatestFunding:        latestFunding,
		LatestFundingAmount:  latestFundingAmount,
		LastRaisedAt:         lastRaisedAt,
		CreatedAt:            createdAt,
		UpdatedAt:            time.Now(),
	}

	esCompany := ElasticCompany{
		Id:               companyUUID,
		Name:             companyName,
		EmployeesCount:   employeesCount,
		Industries:       industries,
		Keywords:         keywords,
		Address:          address,
		AnnualRevenue:    annualRevenue,
		TotalFunding:     totalFunding,
		Technologies:     technologies,
		City:             city,
		State:            state,
		Country:          country,
		LinkedinURL:      linkedinURL,
		Website:          website,
		NormalizedDomain: normalizedDomain,
		CreatedAt:        createdAt,
	}

	return pgCompany, esCompany
}

// ============================================
// Contact Generation
// ============================================

func generateContactRecord(pgCompany PgCompany, esCompany ElasticCompany, companyDomain string) (PgContact, ElasticContact) {
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	email := fmt.Sprintf("%s.%s@%s", toLowerCase(firstName), toLowerCase(lastName), companyDomain)

	// Generate departments (1-3 departments per person)
	numDepartments := rand.Intn(3) + 1
	departments := randomSample(Departments, numDepartments)

	// Generate LinkedIn URL
	linkedinURL := fmt.Sprintf("https://linkedin.com/in/%s-%s-%d",
		toLowerCase(firstName), toLowerCase(lastName), rand.Intn(900000)+100000)

	// Generate created_at date (within last 5 years)
	daysAgo := rand.Intn(1825)
	createdAt := time.Now().AddDate(0, 0, -daysAgo)

	// Generate deterministic UUID based on first_name, last_name, email, linkedin_url
	uuidSeed := fmt.Sprintf("%s%s%s%s", firstName, lastName, email, linkedinURL)
	contactUUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(uuidSeed)).String()

	// Generate extra PG fields
	facebookURL := fmt.Sprintf("https://facebook.com/%s.%s", toLowerCase(firstName), toLowerCase(lastName))
	twitterURL := fmt.Sprintf("https://twitter.com/%s%s", toLowerCase(firstName), toLowerCase(lastName))
	website := fmt.Sprintf("https://%s%s.com", toLowerCase(firstName), toLowerCase(lastName))
	workDirectPhone := gofakeit.Phone()
	homePhone := gofakeit.Phone()
	otherPhone := gofakeit.Phone()
	stage := ContactStages[rand.Intn(len(ContactStages))]

	// Generate contact-specific fields (shared between PG and ES)
	title := Titles[rand.Intn(len(Titles))]
	mobilePhone := gofakeit.Phone()
	emailStatus := EmailStatuses[rand.Intn(len(EmailStatuses))]
	seniority := SeniorityLevels[rand.Intn(len(SeniorityLevels))]

	pgContact := PgContact{
		UUID:            contactUUID,
		FirstName:       firstName,
		LastName:        lastName,
		CompanyID:       pgCompany.UUID,
		Email:           email,
		Title:           title,
		Departments:     departments,
		MobilePhone:     mobilePhone,
		EmailStatus:     emailStatus,
		Seniority:       seniority,
		City:            pgCompany.City,
		State:           pgCompany.State,
		Country:         pgCompany.Country,
		LinkedinURL:     linkedinURL,
		FacebookURL:     facebookURL,
		TwitterURL:      twitterURL,
		Website:         website,
		WorkDirectPhone: workDirectPhone,
		HomePhone:       homePhone,
		OtherPhone:      otherPhone,
		Stage:           stage,
		CreatedAt:       createdAt,
		UpdatedAt:       time.Now(),
	}

	esContact := ElasticContact{
		Id:          contactUUID,
		FirstName:   firstName,
		LastName:    lastName,
		CompanyID:   esCompany.Id,
		Email:       email,
		Title:       title,
		Departments: departments,
		MobilePhone: mobilePhone,
		EmailStatus: emailStatus,
		Seniority:   seniority,
		City:        esCompany.City,
		State:       esCompany.State,
		Country:     esCompany.Country,
		LinkedinURL: linkedinURL,
		CreatedAt:   createdAt,

		// Company details (denormalized from the company for ES searching)
		CompanyName:             esCompany.Name,
		CompanyEmployeesCount:   esCompany.EmployeesCount,
		CompanyIndustries:       esCompany.Industries,
		CompanyKeywords:         esCompany.Keywords,
		CompanyAddress:          esCompany.Address,
		CompanyAnnualRevenue:    esCompany.AnnualRevenue,
		CompanyTotalFunding:     esCompany.TotalFunding,
		CompanyTechnologies:     esCompany.Technologies,
		CompanyCity:             esCompany.City,
		CompanyState:            esCompany.State,
		CompanyCountry:          esCompany.Country,
		CompanyLinkedinURL:      esCompany.LinkedinURL,
		CompanyWebsite:          esCompany.Website,
		CompanyNormalizedDomain: esCompany.NormalizedDomain,
	}

	return pgContact, esContact
}

// ============================================
// Batch Generation
// ============================================

func generateBatch(batchNum int) CompanyBatch {
	pgCompanies := make([]PgCompany, 0, CompanyBatchSize)
	esCompanies := make([]ElasticCompany, 0, CompanyBatchSize)
	pgContacts := make([]PgContact, 0, CompanyBatchSize*ContactsPerCompany)
	esContacts := make([]ElasticContact, 0, CompanyBatchSize*ContactsPerCompany)

	for i := 0; i < CompanyBatchSize; i++ {
		pgCompany, esCompany := generateCompanyRecord()
		pgCompanies = append(pgCompanies, pgCompany)
		esCompanies = append(esCompanies, esCompany)

		// Generate 5 contacts for this company
		for j := 0; j < ContactsPerCompany; j++ {
			domain := strings.TrimPrefix(pgCompany.Website, "https://www.")
			pgContact, esContact := generateContactRecord(pgCompany, esCompany, domain)
			pgContacts = append(pgContacts, pgContact)
			esContacts = append(esContacts, esContact)
		}
	}

	return CompanyBatch{
		BatchNum:    batchNum,
		PgCompanies: pgCompanies,
		PgContacts:  pgContacts,
		EsCompanies: esCompanies,
		EsContacts:  esContacts,
	}
}

// ============================================
// PostgreSQL Operations
// ============================================

func bulkInsertPgCompanies(ctx context.Context, db *bun.DB, companies []PgCompany) (int, int, error) {
	if len(companies) == 0 {
		return 0, 0, nil
	}

	_, err := db.NewInsert().
		Model(&companies).
		On("CONFLICT (uuid) DO NOTHING").
		Exec(ctx)

	if err != nil {
		return 0, len(companies), err
	}

	return len(companies), 0, nil
}

func bulkInsertPgContacts(ctx context.Context, db *bun.DB, contacts []PgContact) (int, int, error) {
	if len(contacts) == 0 {
		return 0, 0, nil
	}

	_, err := db.NewInsert().
		Model(&contacts).
		On("CONFLICT (uuid) DO NOTHING").
		Exec(ctx)

	if err != nil {
		return 0, len(contacts), err
	}

	return len(contacts), 0, nil
}

// ============================================
// Elasticsearch Operations
// ============================================

func bulkIndexEsCompanies(ctx context.Context, es *elasticsearch.Client, companies []ElasticCompany) (int, int, error) {
	if len(companies) == 0 {
		return 0, 0, nil
	}

	var buf strings.Builder

	for _, company := range companies {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": EsCompanyIndex,
				"_id":    company.Id,
			},
		}

		metaJSON, _ := json.Marshal(meta)
		buf.WriteString(string(metaJSON))
		buf.WriteString("\n")

		docJSON, _ := json.Marshal(company)
		buf.WriteString(string(docJSON))
		buf.WriteString("\n")
	}

	req := esapi.BulkRequest{
		Body:    strings.NewReader(buf.String()),
		Refresh: "false",
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		return 0, len(companies), err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, len(companies), fmt.Errorf("bulk request error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, len(companies), err
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return len(companies), 0, nil
	}

	successCount := 0
	failedCount := 0

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		if indexResult, ok := itemMap["index"].(map[string]interface{}); ok {
			if status, ok := indexResult["status"].(float64); ok {
				if status >= 200 && status < 300 {
					successCount++
				} else {
					failedCount++
				}
			}
		}
	}

	return successCount, failedCount, nil
}

func bulkIndexEsContacts(ctx context.Context, es *elasticsearch.Client, contacts []ElasticContact) (int, int, error) {
	if len(contacts) == 0 {
		return 0, 0, nil
	}

	var buf strings.Builder

	for _, contact := range contacts {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": EsContactIndex,
				"_id":    contact.Id,
			},
		}

		metaJSON, _ := json.Marshal(meta)
		buf.WriteString(string(metaJSON))
		buf.WriteString("\n")

		docJSON, _ := json.Marshal(contact)
		buf.WriteString(string(docJSON))
		buf.WriteString("\n")
	}

	req := esapi.BulkRequest{
		Body:    strings.NewReader(buf.String()),
		Refresh: "false",
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		return 0, len(contacts), err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, len(contacts), fmt.Errorf("bulk request error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, len(contacts), err
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return len(contacts), 0, nil
	}

	successCount := 0
	failedCount := 0

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		if indexResult, ok := itemMap["index"].(map[string]interface{}); ok {
			if status, ok := indexResult["status"].(float64); ok {
				if status >= 200 && status < 300 {
					successCount++
				} else {
					failedCount++
				}
			}
		}
	}

	return successCount, failedCount, nil
}

// ============================================
// Producer
// ============================================

func producer(startBatch, endBatch int, batchChan chan<- CompanyBatch, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := startBatch; i < endBatch; i++ {
		batch := generateBatch(i)
		batchChan <- batch
	}
}

// ============================================
// Consumers
// ============================================

func pgConsumer(id int, ctx context.Context, db *bun.DB, batchChan <-chan CompanyBatch, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	for batch := range batchChan {
		// Insert companies
		companySuccess, companyFailed, err := bulkInsertPgCompanies(ctx, db, batch.PgCompanies)
		if err != nil {
			fmt.Printf("PG Consumer %d: Error inserting companies batch %d: %v\n", id, batch.BatchNum+1, err)
		}
		stats.AddPgCompanies(int64(companySuccess), int64(companyFailed))

		// Insert contacts
		contactSuccess, contactFailed, err := bulkInsertPgContacts(ctx, db, batch.PgContacts)
		if err != nil {
			fmt.Printf("PG Consumer %d: Error inserting contacts batch %d: %v\n", id, batch.BatchNum+1, err)
		}
		stats.AddPgContacts(int64(contactSuccess), int64(contactFailed))

		fmt.Printf("PG Consumer %d: Batch %d - Companies: %d, Contacts: %d\n",
			id, batch.BatchNum+1, companySuccess, contactSuccess)
	}
}

func esConsumer(id int, ctx context.Context, es *elasticsearch.Client, batchChan <-chan CompanyBatch, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	for batch := range batchChan {
		// Index companies
		companySuccess, companyFailed, err := bulkIndexEsCompanies(ctx, es, batch.EsCompanies)
		if err != nil {
			fmt.Printf("ES Consumer %d: Error indexing companies batch %d: %v\n", id, batch.BatchNum+1, err)
		}
		stats.AddEsCompanies(int64(companySuccess), int64(companyFailed))

		// Index contacts
		contactSuccess, contactFailed, err := bulkIndexEsContacts(ctx, es, batch.EsContacts)
		if err != nil {
			fmt.Printf("ES Consumer %d: Error indexing contacts batch %d: %v\n", id, batch.BatchNum+1, err)
		}
		stats.AddEsContacts(int64(contactSuccess), int64(contactFailed))

		stats.IncrementBatches()

		fmt.Printf("ES Consumer %d: Batch %d - Companies: %d, Contacts: %d\n",
			id, batch.BatchNum+1, companySuccess, contactSuccess)
	}
}

// ============================================
// Main Function
// ============================================

func main() {
	fmt.Println("==============================================")
	fmt.Println("     Data Generator - Companies & Contacts    ")
	fmt.Println("==============================================")
	fmt.Printf("Total Companies: %d\n", TotalCompanies)
	fmt.Printf("Contacts per Company: %d\n", ContactsPerCompany)
	fmt.Printf("Total Contacts: %d\n", TotalContacts)
	fmt.Printf("Batch Size (Companies): %d\n", CompanyBatchSize)
	fmt.Printf("Total Batches: %d\n", TotalCompanies/CompanyBatchSize)
	fmt.Printf("Producers: %d, PG Consumers: %d, ES Consumers: %d\n\n", NumProducers, NumPgConsumers, NumEsConsumers)

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(time.Now().UnixNano())

	ctx := context.Background()

	// ========================================
	// Connect to PostgreSQL
	// ========================================
	fmt.Println("Connecting to PostgreSQL...")
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(PgDSN),
		pgdriver.WithInsecure(true), // Disable SSL since server doesn't support it
		pgdriver.WithDialTimeout(30*time.Minute),
		pgdriver.WithReadTimeout(30*time.Minute),
		pgdriver.WithWriteTimeout(30*time.Minute),
	))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Error connecting to PostgreSQL: %v\n", err)
		return
	}
	defer db.Close()
	fmt.Println("PostgreSQL connected successfully!")

	// ========================================
	// Connect to Elasticsearch
	// ========================================
	fmt.Println("Connecting to Elasticsearch...")
	esCfg := elasticsearch.Config{
		Addresses: []string{EsURL},
		Username:  EsUsername,
		Password:  EsPassword,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   100,
			ResponseHeaderTimeout: 30 * time.Minute,
			IdleConnTimeout:       30 * time.Minute,
		},
		MaxRetries: 3,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		fmt.Printf("Error creating Elasticsearch client: %v\n", err)
		return
	}

	info, err := es.Info()
	if err != nil {
		fmt.Printf("Error connecting to Elasticsearch: %v\n", err)
		return
	}
	defer info.Body.Close()
	fmt.Println("Elasticsearch connected successfully!")

	for i := 0; i < 8; i++ {
		// ========================================
		// Initialize Statistics
		// ========================================
		totalBatches := TotalCompanies / CompanyBatchSize
		stats := &Stats{
			TotalBatches: int64(totalBatches),
		}

		// ========================================
		// Create Channels
		// ========================================
		// Single channel for generated batches, two channels for PG and ES consumers
		generatedChan := make(chan CompanyBatch, ChannelBufferSize)
		pgChan := make(chan CompanyBatch, ChannelBufferSize)
		esChan := make(chan CompanyBatch, ChannelBufferSize)

		// ========================================
		// Start Timer
		// ========================================
		startTime := time.Now()

		// ========================================
		// Start Producers
		// ========================================
		var producerWg sync.WaitGroup
		batchesPerProducer := totalBatches / NumProducers

		for i := 0; i < NumProducers; i++ {
			producerWg.Add(1)
			startBatch := i * batchesPerProducer
			endBatch := startBatch + batchesPerProducer
			if i == NumProducers-1 {
				endBatch = totalBatches // Last producer handles remaining batches
			}
			go producer(startBatch, endBatch, generatedChan, &producerWg)
		}

		// Close generatedChan when all producers are done
		go func() {
			producerWg.Wait()
			close(generatedChan)
		}()

		// ========================================
		// Fanout: Distribute batches to PG and ES channels
		// ========================================
		go func() {
			for batch := range generatedChan {
				pgChan <- batch
				esChan <- batch
			}
			close(pgChan)
			close(esChan)
		}()

		// ========================================
		// Start Consumers
		// ========================================
		var pgWg sync.WaitGroup
		var esWg sync.WaitGroup

		// Start PostgreSQL consumers
		for i := 0; i < NumPgConsumers; i++ {
			pgWg.Add(1)
			go pgConsumer(i+1, ctx, db, pgChan, stats, &pgWg)
		}

		// Start Elasticsearch consumers
		for i := 0; i < NumEsConsumers; i++ {
			esWg.Add(1)
			go esConsumer(i+1, ctx, es, esChan, stats, &esWg)
		}

		// ========================================
		// Progress Reporter
		// ========================================
		done := make(chan struct{})
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					stats.Print()
				case <-done:
					return
				}
			}
		}()

		// ========================================
		// Wait for Completion
		// ========================================
		pgWg.Wait()
		esWg.Wait()
		close(done)

		// ========================================
		// Final Report
		// ========================================
		elapsed := time.Since(startTime)
		fmt.Println("\n==============================================")
		fmt.Println("              FINAL RESULTS                   ")
		fmt.Println("==============================================")
		fmt.Printf("Time Elapsed: %v\n", elapsed)
		fmt.Printf("Records per second: %.2f\n", float64(TotalCompanies+TotalContacts)/elapsed.Seconds())
		stats.Print()
		fmt.Println("==============================================")
		fmt.Println("              COMPLETED!                      ")
		fmt.Println("==============================================")
	}
}
