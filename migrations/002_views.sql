CREATE MATERIALIZED VIEW IF NOT EXISTS mv_domain_scores AS
SELECT q.domain_id, sr.region_id, sr.district_id, sr.facility_id, sr.period,
SUM(COALESCE(rl.numeric_score,0)) AS numerator,
COUNT(rl.id) FILTER (WHERE rl.is_missing = FALSE) AS denominator,
CASE WHEN COUNT(rl.id) FILTER (WHERE rl.is_missing = FALSE) = 0 THEN 0
ELSE SUM(COALESCE(rl.numeric_score,0)) / COUNT(rl.id) FILTER (WHERE rl.is_missing = FALSE) END AS compliance
FROM responses_long rl
JOIN questions q ON q.id = rl.question_id
JOIN supervision_records sr ON sr.id = rl.record_id
GROUP BY q.domain_id, sr.region_id, sr.district_id, sr.facility_id, sr.period;

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_followup_burden AS
SELECT COALESCE(domain,'Unknown') AS domain, responsibility, status, COUNT(*) AS action_count
FROM follow_up_actions
GROUP BY COALESCE(domain,'Unknown'), responsibility, status;
