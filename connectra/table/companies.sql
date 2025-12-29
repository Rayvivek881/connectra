-- Drop existing table if it exists
DROP TABLE IF EXISTS companies CASCADE;

-- Create table
CREATE TABLE companies
(
    id                      bigserial
        PRIMARY KEY,
    uuid                    text,
    name                    text,
    employees_count         bigint,
    industries              text[],
    keywords                text[],
    address                 text,
    annual_revenue          bigint,
    total_funding           bigint,
    technologies            text[],
    text_search             text,
    created_at              timestamp,
    updated_at              timestamp,
    status                  text DEFAULT '_'::text,
    linkedin_url            text,
    linkedin_sales_url      text,
    facebook_url            text,
    twitter_url             text,
    website                 text,
    company_name_for_emails text,
    phone_number            text,
    latest_funding          text,
    latest_funding_amount   bigint,
    last_raised_at          text,
    city                    text,
    state                   text,
    country                 text
);

-- Set table owner
ALTER TABLE companies
    OWNER TO postgres;

-- Create indexes
CREATE UNIQUE INDEX idx_companies_uuid_unique
    ON companies (uuid);

CREATE INDEX idx_dec_trgm
    ON companies USING gin (text_search gin_trgm_ops);

CREATE INDEX idx_companies_name
    ON companies (name);

CREATE INDEX idx_companies_employees_count
    ON companies (employees_count);

CREATE INDEX idx_companies_annual_revenue
    ON companies (annual_revenue);

CREATE INDEX idx_companies_total_funding
    ON companies (total_funding);

CREATE INDEX idx_companies_industries_gin
    ON companies USING gin (industries);

CREATE INDEX idx_companies_keywords_gin
    ON companies USING gin (keywords);

CREATE INDEX idx_companies_technologies_gin
    ON companies USING gin (technologies);

CREATE INDEX idx_companies_name_trgm
    ON companies USING gin (name gin_trgm_ops);

CREATE INDEX idx_companies_created_at
    ON companies (created_at);

CREATE INDEX idx_companies_annual_revenue_industries
    ON companies (annual_revenue, industries);

CREATE INDEX idx_companies_status
    ON companies (status);

CREATE INDEX idx_companies_status_industries
    ON companies (status, industries);

-- Covering index for uuid-based queries with common columns
-- Note: Removed large TEXT/TEXT[] columns (address, industries, keywords) from INCLUDE
-- to avoid exceeding PostgreSQL's 8191 byte index row size limit.
-- These columns are already indexed separately (GIN indexes for arrays, separate index for name).
-- 
-- If this index still fails due to large company names, use the alternative below:
-- DROP INDEX IF EXISTS idx_companies_uuid_covering;
-- CREATE INDEX IF NOT EXISTS idx_companies_uuid_covering
--     ON companies (uuid)
--     INCLUDE (employees_count, annual_revenue, total_funding);
--
-- Note: Additional optimization indexes are in optimization_indexes.sql
DROP INDEX IF EXISTS idx_companies_uuid_covering;
CREATE INDEX IF NOT EXISTS idx_companies_uuid_covering
    ON companies (uuid)
    INCLUDE (name, employees_count, annual_revenue, total_funding);

-- Indexes for metadata columns (merged from companies_metadata)
-- Performance optimization indexes (from slow query analysis)
CREATE INDEX idx_companies_linkedin_url_btree_exact
    ON companies (linkedin_url)
    WHERE linkedin_url IS NOT NULL;

CREATE INDEX idx_companies_website_not_null
    ON companies (website)
    WHERE website IS NOT NULL AND trim(website) != '';

CREATE INDEX idx_companies_uuid_lookup
    ON companies (uuid)
    WHERE uuid IS NOT NULL;

