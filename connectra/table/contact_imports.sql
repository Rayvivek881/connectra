-- Drop existing tables and types if they exist
DROP TABLE IF EXISTS public.contact_import_errors CASCADE;
DROP TABLE IF EXISTS public.contact_import_jobs CASCADE;
DROP TYPE IF EXISTS import_job_status CASCADE;

-- Create enum type for import job status
CREATE TYPE import_job_status AS ENUM ('pending', 'processing', 'completed', 'failed');

-- Create contact_import_jobs table
CREATE TABLE public.contact_import_jobs
(
    id             bigserial
        PRIMARY KEY,
    job_id         text                NOT NULL,
    file_name      text                NOT NULL,
    file_path      text                NOT NULL,
    total_rows     integer             DEFAULT 0,
    processed_rows integer             DEFAULT 0,
    status         import_job_status  DEFAULT 'pending'::import_job_status,
    error_count    integer             DEFAULT 0,
    message        text,
    created_at     timestamp with time zone NOT NULL DEFAULT now(),
    updated_at     timestamp with time zone NOT NULL DEFAULT now(),
    completed_at   timestamp with time zone
);

-- Set table owner
ALTER TABLE public.contact_import_jobs
    OWNER TO postgres;

-- Create indexes for contact_import_jobs
CREATE UNIQUE INDEX idx_contact_import_jobs_job_id_unique
    ON public.contact_import_jobs (job_id);

CREATE INDEX idx_contact_import_jobs_job_id
    ON public.contact_import_jobs (job_id);

CREATE INDEX idx_contact_import_jobs_status
    ON public.contact_import_jobs (status);

CREATE INDEX idx_contact_import_jobs_created_at
    ON public.contact_import_jobs (created_at);

-- Create contact_import_errors table
CREATE TABLE public.contact_import_errors
(
    id            bigserial
        PRIMARY KEY,
    job_id        bigint              NOT NULL,
    row_number    integer             NOT NULL,
    error_message text                NOT NULL,
    payload       text,
    created_at    timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT fk_contact_import_errors_job_id
        FOREIGN KEY (job_id)
            REFERENCES public.contact_import_jobs (id)
            ON DELETE CASCADE
);

-- Set table owner
ALTER TABLE public.contact_import_errors
    OWNER TO postgres;

-- Create indexes for contact_import_errors
CREATE INDEX idx_contact_import_errors_job_id
    ON public.contact_import_errors (job_id);

CREATE INDEX idx_contact_import_errors_created_at
    ON public.contact_import_errors (created_at);

