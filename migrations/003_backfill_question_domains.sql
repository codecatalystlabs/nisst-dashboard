-- Backfill questions.domain_id from common NISST / ODK column naming (prefixes).
-- Run after 001_init.sql. Safe to run multiple times (only updates NULL domain_id).

-- D1 Leadership and Governance
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D1' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd1_') = 1 OR strpos(lower(q.source_column), '/d1_') > 0
    OR strpos(lower(q.source_column), 'lg_') = 1 OR strpos(lower(q.source_column), '/lg_') > 0
    OR strpos(lower(q.source_column), 'ldr_') = 1 OR strpos(lower(q.source_column), '/ldr_') > 0
  );

-- D2 Human Resources for Health (hrh_ before hr_: otherwise hr_ matches hrh_*)
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D2' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd2_') = 1 OR strpos(lower(q.source_column), '/d2_') > 0
    OR strpos(lower(q.source_column), 'hrh_') = 1 OR strpos(lower(q.source_column), '/hrh_') > 0
    OR (
      (strpos(lower(q.source_column), 'hr_') = 1 OR strpos(lower(q.source_column), '/hr_') > 0)
      AND NOT (strpos(lower(q.source_column), 'hrh_') = 1 OR strpos(lower(q.source_column), '/hrh_') > 0)
    )
  );

-- D3 Medicines and Health Supplies
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D3' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd3_') = 1 OR strpos(lower(q.source_column), '/d3_') > 0
    OR strpos(lower(q.source_column), 'mhs_') = 1 OR strpos(lower(q.source_column), '/mhs_') > 0
    OR strpos(lower(q.source_column), 'med_') = 1 OR strpos(lower(q.source_column), '/med_') > 0
    OR strpos(lower(q.source_column), 'ms_') = 1 OR strpos(lower(q.source_column), '/ms_') > 0
  );

-- D4 Health Financing
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D4' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd4_') = 1 OR strpos(lower(q.source_column), '/d4_') > 0
    OR strpos(lower(q.source_column), 'hf_') = 1 OR strpos(lower(q.source_column), '/hf_') > 0
    OR strpos(lower(q.source_column), 'fin_') = 1 OR strpos(lower(q.source_column), '/fin_') > 0
  );

-- D5 Health Information Management
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D5' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd5_') = 1 OR strpos(lower(q.source_column), '/d5_') > 0
    OR strpos(lower(q.source_column), 'him_') = 1 OR strpos(lower(q.source_column), '/him_') > 0
    OR strpos(lower(q.source_column), 'hin_') = 1 OR strpos(lower(q.source_column), '/hin_') > 0
  );

-- D6 Health Infrastructure
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D6' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd6_') = 1 OR strpos(lower(q.source_column), '/d6_') > 0
    OR strpos(lower(q.source_column), 'infra_') = 1 OR strpos(lower(q.source_column), '/infra_') > 0
    OR strpos(lower(q.source_column), 'inf_') = 1 OR strpos(lower(q.source_column), '/inf_') > 0
  );

-- D7 Service Delivery
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D7' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd7_') = 1 OR strpos(lower(q.source_column), '/d7_') > 0
    OR strpos(lower(q.source_column), 'srv_') = 1 OR strpos(lower(q.source_column), '/srv_') > 0
    OR strpos(lower(q.source_column), 'serv_') = 1 OR strpos(lower(q.source_column), '/serv_') > 0
    OR strpos(lower(q.source_column), 'sd_') = 1 OR strpos(lower(q.source_column), '/sd_') > 0
  );

-- D8 Quality of Care and Safety
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D8' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd8_') = 1 OR strpos(lower(q.source_column), '/d8_') > 0
    OR strpos(lower(q.source_column), 'qoc_') = 1 OR strpos(lower(q.source_column), '/qoc_') > 0
    OR strpos(lower(q.source_column), 'qos_') = 1 OR strpos(lower(q.source_column), '/qos_') > 0
    OR strpos(lower(q.source_column), 'qc_') = 1 OR strpos(lower(q.source_column), '/qc_') > 0
  );

-- D9 Follow-up Action Areas (rare on main form)
UPDATE questions q
SET domain_id = (SELECT id FROM domains WHERE code = 'D9' LIMIT 1)
WHERE q.domain_id IS NULL
  AND (
    strpos(lower(q.source_column), 'd9_') = 1 OR strpos(lower(q.source_column), '/d9_') > 0
    OR strpos(lower(q.source_column), 'fua_') = 1 OR strpos(lower(q.source_column), '/fua_') > 0
    OR strpos(lower(q.source_column), 'fu_') = 1 OR strpos(lower(q.source_column), '/fu_') > 0
  );
