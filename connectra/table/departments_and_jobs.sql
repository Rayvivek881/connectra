-- Drop existing table if it exists
DROP TABLE IF EXISTS public.departments_and_jobs CASCADE;

-- Create table
CREATE TABLE public.departments_and_jobs
(
    id           bigserial
        PRIMARY KEY,
    department   text DEFAULT '_'::text,
    job_function text DEFAULT '_'::text,
    uuid         varchar(50) NOT NULL
);

-- Set table owner
ALTER TABLE public.departments_and_jobs
    OWNER TO postgres;

