-- Drop existing table if it exists
DROP TABLE IF EXISTS public.contacts CASCADE;

-- Create table
CREATE TABLE public.contacts
(
    id                bigserial
        PRIMARY KEY,
    uuid              text,
    first_name        text,
    last_name         text,
    company_id        text,
    email             text,
    title             text,
    departments       text[],
    mobile_phone      text,
    email_status      text,
    text_search       text,
    created_at        timestamp,
    updated_at        timestamp,
    seniority         text DEFAULT '_'::text,
    status            text DEFAULT '_'::text,
    linkedin_url      text DEFAULT '_'::text,
    linkedin_sales_url text DEFAULT '_'::text,
    facebook_url      text DEFAULT '_'::text,
    twitter_url       text DEFAULT '_'::text,
    website           text DEFAULT '_'::text,
    work_direct_phone text DEFAULT '_'::text,
    home_phone        text DEFAULT '_'::text,
    city              text DEFAULT '_'::text,
    state             text DEFAULT '_'::text,
    country           text DEFAULT '_'::text,
    other_phone       text DEFAULT '_'::text,
    stage             text DEFAULT '_'::text
);

-- Set table owner
ALTER TABLE public.contacts
    OWNER TO postgres;

-- Create indexes
CREATE UNIQUE INDEX idx_contacts_uuid_unique
    ON public.contacts (uuid);

CREATE INDEX idx_contacts_first_name
    ON public.contacts (first_name);

CREATE INDEX idx_contacts_last_name
    ON public.contacts (last_name);

CREATE INDEX idx_contacts_company_id
    ON public.contacts (company_id);

CREATE INDEX idx_contacts_email
    ON public.contacts (email);

CREATE INDEX idx_contacts_mobile_phone
    ON public.contacts (mobile_phone);

CREATE INDEX idx_contacts_email_status
    ON public.contacts (email_status);

CREATE INDEX idx_contacts_title
    ON public.contacts (title);

CREATE INDEX idx_contacts_title_trgm
    ON public.contacts USING gin (title public.gin_trgm_ops);

CREATE INDEX idx_contacts_email_company
    ON public.contacts (email, company_id);

CREATE INDEX idx_contacts_name_company
    ON public.contacts (first_name, last_name, company_id);

CREATE INDEX idx_contacts_created_at
    ON public.contacts (created_at);

CREATE INDEX idx_contacts_seniority
    ON public.contacts (seniority);

CREATE INDEX idx_contacts_seniority_company_id
    ON public.contacts (seniority, company_id);

CREATE INDEX idx_contacts_departments_gin
    ON public.contacts USING gin (departments);

CREATE INDEX idx_contacts_company_department
    ON public.contacts (company_id, departments);

CREATE INDEX idx_contacts_seniority_department
    ON public.contacts (seniority, departments);

CREATE INDEX idx_contacts_dec_trgm
    ON public.contacts USING gin (text_search public.gin_trgm_ops);

CREATE INDEX idx_contacts_status
    ON public.contacts (status);

CREATE INDEX idx_contacts_status_company_id
    ON public.contacts (status, company_id);

-- Indexes for metadata columns (merged from contacts_metadata)
-- B-tree index for linkedin_url exact/prefix matches
-- Performance optimization indexes (from slow query analysis)
CREATE INDEX idx_contacts_linkedin_url_btree
    ON public.contacts (linkedin_url)
    WHERE linkedin_url IS NOT NULL AND linkedin_url != '_';

CREATE INDEX idx_contacts_uuid_lookup
    ON public.contacts (uuid)
    WHERE uuid IS NOT NULL;

