# Configuration Analysis

## Analysis of Provided Configuration

### Issues Identified

1. **Duplicates Found:**
   - `APP_ENV` appears twice (both set to `development`)
   - `RUN_COMMAND` appears twice (once as `api-server`, once as `jobs first_time`)
   - All other configuration values appear twice

2. **Typo in Codebase:**
   - `S3_UPLOAD_FILE_PATH_PRIFIX` should be `PREFIX` but codebase uses `PRIFIX` consistently
   - Keeping as-is to maintain compatibility with existing code

3. **Missing Configuration:**
   - `RUN_COMMAND` is not in Viper config (used only in Dockerfile as environment variable)
   - This is correct - it's a Docker runtime variable, not application config

4. **Configuration Structure:**
   - Two separate configs provided (API server and Jobs)
   - Both use same database/Elasticsearch/S3 settings
   - Only difference is `RUN_COMMAND` and some debug flags

### Configuration Mapping

| Provided Config | Viper Config | Status |
|----------------|--------------|--------|
| `APP_ENV` | `APP_ENV` | ✅ Matches |
| `API_KEY` | `API_KEY` | ✅ Matches |
| `MAX_REQUESTS_PER_MINUTE` | `MAX_REQUESTS_PER_MINUTE` | ✅ Matches |
| `MEMORY_LOG_INTERVAL_SECONDS` | `MEMORY_LOG_INTERVAL_SECONDS` | ✅ Matches |
| `PG_DB_CONNECTION` | `PG_DB_CONNECTION` | ✅ Matches |
| `PG_DB_HOST` | `PG_DB_HOST` | ✅ Matches |
| `PG_DB_PORT` | `PG_DB_PORT` | ✅ Matches |
| `PG_DB_DATABASE` | `PG_DB_DATABASE` | ✅ Matches |
| `PG_DB_USERNAME` | `PG_DB_USERNAME` | ✅ Matches |
| `PG_DB_PASSWORD` | `PG_DB_PASSWORD` | ✅ Matches |
| `PG_DB_DEBUG` | `PG_DB_DEBUG` | ✅ Matches |
| `PG_DB_SSL` | `PG_DB_SSL` | ✅ Matches |
| `ELASTICSEARCH_CONNECTION` | `ELASTICSEARCH_CONNECTION` | ✅ Matches |
| `ELASTICSEARCH_HOST` | `ELASTICSEARCH_HOST` | ✅ Matches |
| `ELASTICSEARCH_PORT` | `ELASTICSEARCH_PORT` | ✅ Matches |
| `ELASTICSEARCH_USERNAME` | `ELASTICSEARCH_USERNAME` | ✅ Matches |
| `ELASTICSEARCH_PASSWORD` | `ELASTICSEARCH_PASSWORD` | ✅ Matches |
| `ELASTICSEARCH_DEBUG` | `ELASTICSEARCH_DEBUG` | ✅ Matches |
| `ELASTICSEARCH_SSL` | `ELASTICSEARCH_SSL` | ✅ Matches |
| `ELASTICSEARCH_AUTH` | `ELASTICSEARCH_AUTH` | ✅ Matches |
| `S3_ACCESS_KEY` | `S3_ACCESS_KEY` | ✅ Matches |
| `S3_SECRET_KEY` | `S3_SECRET_KEY` | ✅ Matches |
| `S3_REGION` | `S3_REGION` | ✅ Matches |
| `S3_BUCKET` | `S3_BUCKET` | ✅ Matches |
| `S3_ENDPOINT` | `S3_ENDPOINT` | ✅ Matches |
| `S3_SSL` | `S3_SSL` | ✅ Matches |
| `S3_DEBUG` | `S3_DEBUG` | ✅ Matches |
| `S3_UPLOAD_URL_TTL_HOURS` | `S3_UPLOAD_URL_TTL_HOURS` | ✅ Matches |
| `S3_UPLOAD_FILE_PATH_PRIFIX` | `S3_UPLOAD_FILE_PATH_PRIFIX` | ✅ Matches (typo in codebase) |
| `JOB_IN_QUEUE_SIZE` | `JOB_IN_QUEUE_SIZE` | ✅ Matches |
| `PARALLEL_JOBS` | `PARALLEL_JOBS` | ✅ Matches |
| `TICKER_INTERVAL_MINUTES` | `TICKER_INTERVAL_MINUTES` | ✅ Matches |
| `BATCH_SIZE_FOR_INSERTION` | `BATCH_SIZE_FOR_INSERTION` | ✅ Matches |
| `JOB_TYPE` | `JOB_TYPE` | ✅ Matches |
| `RUN_COMMAND` | N/A | ⚠️ Docker-only variable |

### Production Configuration Values

- **Database**: Production PostgreSQL at `98.81.200.121:5432`
- **Elasticsearch**: Production instance at `54.198.202.181:9200`
- **S3**: AWS S3 bucket `tkdrawcsvdata` in `us-east-1`
- **API Key**: `3e6b8811-40c2-46e7-8d7c-e7e038e86071`
- **Rate Limit**: 180 requests per minute
- **Job Workers**: 5 parallel jobs
- **Batch Size**: 100 records per batch

### Recommendations

1. ✅ Use single `.env` file (RUN_COMMAND set via Docker)
2. ✅ Keep `PRIFIX` typo for compatibility
3. ✅ Remove duplicate entries
4. ✅ Organize by sections (App, Database, Elasticsearch, S3, Jobs)
5. ✅ Keep local configs commented for easy switching
