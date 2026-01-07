# Configuration Setup Guide

## Overview

This guide explains how to set up the environment configuration for the Connectra API.

## Configuration Files

- **`.example.env`** - Template with example values (safe to commit)
- **`.env.production`** - Production configuration template (contains real credentials)
- **`.env`** - Actual configuration file (gitignored, created from templates)

## Quick Setup

### For Production

```bash
# Copy production template to .env
cp .env.production .env
```

### For Local Development

```bash
# Copy example template to .env
cp .example.env .env
# Edit .env with your local database/Elasticsearch credentials
```

## Configuration Structure

### Application Configuration

```env
APP_ENV=development
API_KEY=your-secret-api-key
MAX_REQUESTS_PER_MINUTE=180
MEMORY_LOG_INTERVAL_SECONDS=900
```

### PostgreSQL Database

```env
PG_DB_CONNECTION=postgres://user:password@host:port/database
PG_DB_HOST=host
PG_DB_PORT=5432
PG_DB_DATABASE=database
PG_DB_USERNAME=user
PG_DB_PASSWORD=password
PG_DB_DEBUG=false
PG_DB_SSL=false
```

### Elasticsearch

```env
ELASTICSEARCH_CONNECTION=http://host:9200
ELASTICSEARCH_HOST=host
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=password
ELASTICSEARCH_DEBUG=false
ELASTICSEARCH_SSL=false
ELASTICSEARCH_AUTH=true
```

### S3 Storage

```env
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_REGION=us-east-1
S3_BUCKET=your-bucket-name
S3_ENDPOINT=
S3_SSL=true
S3_DEBUG=false
S3_UPLOAD_URL_TTL_HOURS=24
S3_UPLOAD_FILE_PATH_PRIFIX=connectra/
```

**Note**: `S3_UPLOAD_FILE_PATH_PRIFIX` has a typo (should be PREFIX) but is kept as-is for compatibility with the codebase.

### Job Configuration

```env
JOB_IN_QUEUE_SIZE=100
PARALLEL_JOBS=5
TICKER_INTERVAL_MINUTES=1
BATCH_SIZE_FOR_INSERTION=100
JOB_TYPE=normal
```

## Docker Runtime Configuration

`RUN_COMMAND` is **not** in the `.env` file. It's set as a Docker environment variable:

```bash
# For API server
docker run -e RUN_COMMAND=api-server ...

# For job processor (first time)
docker run -e RUN_COMMAND="jobs first_time" ...

# For job processor (retry)
docker run -e RUN_COMMAND="jobs retry" ...
```

## Verification

After creating `.env`, verify configuration loads correctly:

```bash
# Build the application (will load .env via Viper)
go build ./...

# Or run the application
go run main.go api-server
```

## Security Notes

⚠️ **Important Security Considerations:**

1. **Never commit `.env` file** - It's in `.gitignore` for a reason
2. **Rotate credentials regularly** - Especially API keys and database passwords
3. **Use environment-specific configs** - Different values for dev/staging/prod
4. **Restrict file permissions** - `chmod 600 .env` on Unix systems
5. **Use secrets management** - Consider AWS Secrets Manager, HashiCorp Vault, etc. for production

## Troubleshooting

### Configuration Not Loading

- Ensure `.env` file exists in the project root
- Check file permissions
- Verify no syntax errors (missing quotes, etc.)
- Check Viper logs for loading errors

### Wrong Values Being Used

- Viper loads environment variables with higher priority than `.env`
- Check if environment variables are set: `env | grep PG_DB`
- Unset conflicting environment variables if needed

### Typo in S3_UPLOAD_FILE_PATH_PRIFIX

This is intentional - the codebase uses `PRIFIX` (typo) consistently. Changing it would require updating:
- `conf/viper.go`
- All documentation
- Existing deployments

## Production Configuration

The production configuration is stored in `.env.production` with:
- Production database: `98.81.200.121:5432`
- Production Elasticsearch: `54.198.202.181:9200`
- AWS S3 bucket: `tkdrawcsvdata`
- API Key: `3e6b8811-40c2-46e7-8d7c-e7e038e86071`

**⚠️ Keep this file secure and never commit it to version control!**
