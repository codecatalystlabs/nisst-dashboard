CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS upload_batches (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  batch_code TEXT UNIQUE NOT NULL,
  file_type TEXT NOT NULL CHECK (file_type IN ('main','followup','metadata')),
  file_name TEXT NOT NULL,
  file_hash TEXT NOT NULL,
  uploader TEXT NOT NULL,
  dry_run BOOLEAN NOT NULL DEFAULT FALSE,
  status TEXT NOT NULL DEFAULT 'pending',
  total_rows INT NOT NULL DEFAULT 0,
  duplicate_rows INT NOT NULL DEFAULT 0,
  error_rows INT NOT NULL DEFAULT 0,
  imported_rows INT NOT NULL DEFAULT 0,
  import_summary JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  processed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS regions (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name TEXT UNIQUE NOT NULL);
CREATE TABLE IF NOT EXISTS districts (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), region_id UUID REFERENCES regions(id), name TEXT NOT NULL, UNIQUE(region_id, name));
CREATE TABLE IF NOT EXISTS facilities (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), region_id UUID REFERENCES regions(id), district_id UUID REFERENCES districts(id), name TEXT NOT NULL, level TEXT, UNIQUE(region_id,district_id,name));

CREATE TABLE IF NOT EXISTS domains (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT UNIQUE NOT NULL, name TEXT NOT NULL, display_order INT NOT NULL DEFAULT 0);
CREATE TABLE IF NOT EXISTS subdomains (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), domain_id UUID REFERENCES domains(id), code TEXT NOT NULL, name TEXT NOT NULL, display_order INT NOT NULL DEFAULT 0, UNIQUE(domain_id, code));
CREATE TABLE IF NOT EXISTS questions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  domain_id UUID REFERENCES domains(id),
  subdomain_id UUID REFERENCES subdomains(id),
  source_column TEXT UNIQUE NOT NULL,
  label TEXT NOT NULL,
  question_type TEXT NOT NULL,
  option_list TEXT,
  scoreable BOOLEAN NOT NULL DEFAULT FALSE,
  weight NUMERIC(10,4) NOT NULL DEFAULT 1.0,
  yes_score NUMERIC(10,4) NOT NULL DEFAULT 1.0,
  no_score NUMERIC(10,4) NOT NULL DEFAULT 0.0,
  include_null_in_denominator BOOLEAN NOT NULL DEFAULT FALSE,
  display_order INT NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS question_options (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), question_id UUID REFERENCES questions(id), option_value TEXT NOT NULL, option_label TEXT NOT NULL, option_score NUMERIC(10,4), display_order INT NOT NULL DEFAULT 0);

CREATE TABLE IF NOT EXISTS supervision_records (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  upload_batch_id UUID REFERENCES upload_batches(id),
  source_key TEXT UNIQUE NOT NULL,
  submission_date TIMESTAMPTZ,
  period TEXT,
  region_id UUID REFERENCES regions(id),
  district_id UUID REFERENCES districts(id),
  facility_id UUID REFERENCES facilities(id),
  level TEXT,
  raw_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS follow_up_actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  upload_batch_id UUID REFERENCES upload_batches(id),
  source_key TEXT,
  parent_key TEXT NOT NULL,
  linked_record_id UUID REFERENCES supervision_records(id),
  domain TEXT,
  challenge TEXT,
  intervention TEXT,
  responsibility TEXT,
  timelines TEXT,
  status TEXT DEFAULT 'open',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS responses_long (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  record_id UUID REFERENCES supervision_records(id),
  question_id UUID REFERENCES questions(id),
  raw_value TEXT,
  normalized_value TEXT,
  numeric_score NUMERIC(10,4),
  is_missing BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(record_id, question_id)
);
CREATE TABLE IF NOT EXISTS comments (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), record_id UUID REFERENCES supervision_records(id), question_id UUID REFERENCES questions(id), domain_id UUID REFERENCES domains(id), text_value TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE IF NOT EXISTS import_errors (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), upload_batch_id UUID REFERENCES upload_batches(id), row_number INT NOT NULL, source_column TEXT, error_code TEXT NOT NULL, error_detail TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());

CREATE TABLE IF NOT EXISTS indicator_definitions (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT UNIQUE NOT NULL, name TEXT NOT NULL, expression TEXT NOT NULL, entity_scope TEXT NOT NULL, enabled BOOLEAN NOT NULL DEFAULT TRUE);
CREATE TABLE IF NOT EXISTS indicator_results (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), indicator_id UUID REFERENCES indicator_definitions(id), entity_type TEXT NOT NULL, entity_id TEXT NOT NULL, period TEXT, value NUMERIC(14,4) NOT NULL, metadata JSONB NOT NULL DEFAULT '{}'::jsonb, computed_at TIMESTAMPTZ NOT NULL DEFAULT NOW());

CREATE INDEX IF NOT EXISTS idx_records_filters ON supervision_records(region_id,district_id,facility_id,period,submission_date);
CREATE INDEX IF NOT EXISTS idx_followups_parent ON follow_up_actions(parent_key,linked_record_id);
