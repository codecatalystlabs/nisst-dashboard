# NISST Analytics Platform

Full-stack analytics foundation for the National Integrated Supportive Supervision Tool.

## Deliverables included
- PostgreSQL schema + analytics materialized views (`migrations/`)
- Metadata seed scripts (`seeds/`)
- Backend clean architecture scaffold (`backend/`)
- Frontend Next.js App Router executive dashboard scaffold (`frontend/`)
- Dual upload fixtures (`fixtures/`)
- API route map and Swagger starter spec (`backend/docs/swagger.yaml`)

## Backend architecture
- `cmd/server` bootstrap
- `internal/config` env loading
- `internal/database` GORM postgres connector
- `internal/repository` persistence registry
- `internal/service` import + scoring + metadata + analytics services
- `internal/handler` REST handlers for upload/metadata/records/analytics/dashboard
- `internal/middleware` centralized error handling and structured request logging
- `internal/importer` normalization helpers for comment/score value classification
- `internal/jobs` background reprocess job interface

## Key ingestion design
1. Main and follow-up uploads are separate pipelines.
2. Follow-ups are linked by `PARENT_KEY -> meta-instanceID` using deferred linking support.
3. CSV header validation is enforced before import.
4. Dry-run mode is supported via query parameter.
5. Raw row payload is preserved in `supervision_records.raw_payload` for traceability.
6. Normalized analytics flow targets `responses_long`, `comments`, and `follow_up_actions`.

## API groups
- Uploads: `/api/v1/uploads/*`
- Metadata: `/api/v1/metadata/*`
- Records: `/api/v1/records*`
- Analytics: `/api/v1/analytics/*`
- Dashboard: `/api/v1/dashboard/*`
- Swagger UI: `/swagger`

## Frontend UX
- Premium Midas theme (navy/white/gold)
- Sticky executive navigation
- Overview KPIs + charts + performance table
- Domain pages and service-delivery sub-navigation
- Dedicated uploads and metadata pages

## Start locally
1. Copy `.env.example` to `.env`
2. `docker compose up --build`
3. Backend: `http://localhost:8080/swagger`
4. Frontend: `http://localhost:3000`

## Testing
- Unit tests for scoring and import validation in `backend/internal/service`

## Next implementation steps
- Persist upload batches and import errors through repositories.
- Implement XLSForm parser to auto-populate questions/options/domains.
- Implement normalization transaction pipeline and idempotent upsert strategy.
- Add integration tests for `/uploads/main` and `/uploads/followups`.
- Bind frontend filters/charts to live analytics endpoints.
