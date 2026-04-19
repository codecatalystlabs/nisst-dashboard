package service

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"nisst/internal/repository"
)

type ImportService struct{ repos *repository.Registry }
type AnalyticsService struct{ repos *repository.Registry }

type ImportSummary struct {
	BatchID       string `json:"batch_id"`
	TotalRows     int    `json:"total_rows"`
	ImportedRows  int    `json:"imported_rows"`
	DuplicateRows int    `json:"duplicate_rows"`
	ErrorRows     int    `json:"error_rows"`
	FileHash      string `json:"file_hash"`
	Status        string `json:"status"`
	DryRun        bool   `json:"dry_run"`
}

type ImportInput struct {
	FileType string
	FileName string
	Uploader string
	DryRun   bool
	Reader   io.Reader
}

type UploadBatchDTO struct {
	ID            string     `json:"id"`
	BatchCode     string     `json:"batch_code"`
	FileType      string     `json:"file_type"`
	FileName      string     `json:"file_name"`
	FileHash      string     `json:"file_hash"`
	Uploader      string     `json:"uploader"`
	DryRun        bool       `json:"dry_run"`
	Status        string     `json:"status"`
	TotalRows     int        `json:"total_rows"`
	ImportedRows  int        `json:"imported_rows"`
	DuplicateRows int        `json:"duplicate_rows"`
	ErrorRows     int        `json:"error_rows"`
	CreatedAt     time.Time  `json:"created_at"`
	ProcessedAt   *time.Time `json:"processed_at"`
}

type ImportErrorDTO struct {
	ID           string    `json:"id"`
	RowNumber    int       `json:"row_number"`
	SourceColumn string    `json:"source_column"`
	ErrorCode    string    `json:"error_code"`
	ErrorDetail  string    `json:"error_detail"`
	CreatedAt    time.Time `json:"created_at"`
}

var mainHeaders = []string{"SubmissionDate", "facility_name", "level", "region", "district", "period", "meta-instanceID"}
var followHeaders = []string{"domain", "challenge", "intervention", "responsibility", "timelines", "PARENT_KEY", "KEY"}

func NewImportService(r *repository.Registry) *ImportService     { return &ImportService{repos: r} }
func NewMetadataService(r *repository.Registry) *MetadataService { return &MetadataService{repos: r} }
func NewAnalyticsService(r *repository.Registry) *AnalyticsService {
	return &AnalyticsService{repos: r}
}
func NewRecordService(r *repository.Registry) *RecordService { return &RecordService{repos: r} }

type MetadataService struct{ repos *repository.Registry }
type RecordService struct{ repos *repository.Registry }

type uploadBatchModel struct {
	ID            string     `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	BatchCode     string     `gorm:"column:batch_code"`
	FileType      string     `gorm:"column:file_type"`
	FileName      string     `gorm:"column:file_name"`
	FileHash      string     `gorm:"column:file_hash"`
	Uploader      string     `gorm:"column:uploader"`
	DryRun        bool       `gorm:"column:dry_run"`
	Status        string     `gorm:"column:status"`
	TotalRows     int        `gorm:"column:total_rows"`
	ImportedRows  int        `gorm:"column:imported_rows"`
	DuplicateRows int        `gorm:"column:duplicate_rows"`
	ErrorRows     int        `gorm:"column:error_rows"`
	ImportSummary []byte     `gorm:"column:import_summary"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	ProcessedAt   *time.Time `gorm:"column:processed_at"`
}

func (uploadBatchModel) TableName() string { return "upload_batches" }

type importErrorModel struct {
	ID            string    `gorm:"column:id"`
	UploadBatchID string    `gorm:"column:upload_batch_id"`
	RowNumber     int       `gorm:"column:row_number"`
	SourceColumn  string    `gorm:"column:source_column"`
	ErrorCode     string    `gorm:"column:error_code"`
	ErrorDetail   string    `gorm:"column:error_detail"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (importErrorModel) TableName() string { return "import_errors" }

func (s *ImportService) ProcessUpload(ctx context.Context, in ImportInput) (ImportSummary, error) {
	content, err := io.ReadAll(in.Reader)
	if err != nil {
		return ImportSummary{}, err
	}
	h := sha256.Sum256(content)
	fileHash := hex.EncodeToString(h[:])

	var existing uploadBatchModel
	err = s.repos.DB.WithContext(ctx).Where("file_type = ? AND file_hash = ? AND dry_run = ?", in.FileType, fileHash, in.DryRun).Order("created_at DESC").First(&existing).Error
	if err == nil {
		return ImportSummary{
			BatchID:       existing.ID,
			TotalRows:     existing.TotalRows,
			ImportedRows:  existing.ImportedRows,
			DuplicateRows: existing.DuplicateRows,
			ErrorRows:     existing.ErrorRows,
			FileHash:      existing.FileHash,
			Status:        "idempotent_skip",
			DryRun:        existing.DryRun,
		}, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ImportSummary{}, err
	}

	reader := csv.NewReader(strings.NewReader(string(content)))
	headers, err := reader.Read()
	if err != nil {
		return ImportSummary{}, err
	}
	if err := validateHeaders(in.FileType, headers); err != nil {
		return ImportSummary{}, err
	}

	batchID := uuid.NewString()
	batchCode := "batch_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	summary := ImportSummary{
		BatchID:  batchID,
		FileHash: fileHash,
		Status:   "completed",
		DryRun:   in.DryRun,
	}

	err = s.repos.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		batch := uploadBatchModel{
			ID:        batchID,
			BatchCode: batchCode,
			FileType:  in.FileType,
			FileName:  in.FileName,
			FileHash:  fileHash,
			Uploader:  in.Uploader,
			DryRun:    in.DryRun,
			Status:    "processing",
		}
		if err := tx.Create(&batch).Error; err != nil {
			return err
		}

		rows, rowErr := reader.ReadAll()
		if rowErr != nil {
			return rowErr
		}
		summary.TotalRows = len(rows)

		for i, row := range rows {
			rowMap := toRowMap(headers, row)
			switch in.FileType {
			case "main":
				if err := s.processMainRow(tx, in.DryRun, batchID, rowMap); err != nil {
					if strings.Contains(err.Error(), "duplicate") {
						summary.DuplicateRows++
					} else {
						summary.ErrorRows++
					}
					if e := s.logImportError(tx, batchID, i+2, "", "MAIN_ROW_ERROR", err.Error()); e != nil {
						return e
					}
					continue
				}
			case "followup":
				if err := s.processFollowupRow(tx, in.DryRun, batchID, rowMap); err != nil {
					if strings.Contains(err.Error(), "duplicate") {
						summary.DuplicateRows++
					} else {
						summary.ErrorRows++
					}
					if e := s.logImportError(tx, batchID, i+2, "", "FOLLOWUP_ROW_ERROR", err.Error()); e != nil {
						return e
					}
					continue
				}
			default:
				return fmt.Errorf("unsupported file type: %s", in.FileType)
			}
			summary.ImportedRows++
		}

		if in.FileType == "main" && !in.DryRun {
			if err := tx.Exec(`
				UPDATE follow_up_actions fu
				SET linked_record_id = sr.id
				FROM supervision_records sr
				WHERE fu.linked_record_id IS NULL
				  AND fu.parent_key = sr.source_key
			`).Error; err != nil {
				return err
			}
		}

		payload := fmt.Sprintf(`{"total_rows":%d,"imported_rows":%d,"duplicate_rows":%d,"error_rows":%d}`, summary.TotalRows, summary.ImportedRows, summary.DuplicateRows, summary.ErrorRows)
		now := time.Now().UTC()
		if err := tx.Model(&uploadBatchModel{}).
			Where("id = ?", batchID).
			Updates(map[string]any{
				"status":         "completed",
				"total_rows":     summary.TotalRows,
				"imported_rows":  summary.ImportedRows,
				"duplicate_rows": summary.DuplicateRows,
				"error_rows":     summary.ErrorRows,
				"import_summary": []byte(payload),
				"processed_at":   &now,
			}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return ImportSummary{}, err
	}

	return summary, nil
}

func (s *ImportService) processMainRow(tx *gorm.DB, dryRun bool, batchID string, row map[string]string) error {
	sourceKey := strings.TrimSpace(row["meta-instanceID"])
	if sourceKey == "" {
		return errors.New("meta-instanceID is required")
	}

	var existing int64
	if err := tx.Table("supervision_records").Where("source_key = ?", sourceKey).Count(&existing).Error; err != nil {
		return err
	}
	if existing > 0 {
		return errors.New("duplicate source key")
	}
	if dryRun {
		return nil
	}

	regionID, err := s.ensureRegion(tx, row["region"])
	if err != nil {
		return err
	}
	districtID, err := s.ensureDistrict(tx, regionID, row["district"])
	if err != nil {
		return err
	}
	facilityID, err := s.ensureFacility(tx, regionID, districtID, row["facility_name"], row["level"])
	if err != nil {
		return err
	}

	var recordID string
	if err := tx.Raw(`
		INSERT INTO supervision_records (
		  id, upload_batch_id, source_key, submission_date, period, region_id, district_id, facility_id, level, raw_payload
		) VALUES (
		  gen_random_uuid(), ?, ?, NULLIF(?, '')::timestamptz, ?, ?, ?, ?, ?, ?::jsonb
		)
		RETURNING id
	`, batchID, sourceKey, row["SubmissionDate"], row["period"], regionID, districtID, facilityID, row["level"], toJSON(row)).Scan(&recordID).Error; err != nil {
		return err
	}
	return s.normalizeMainRow(tx, recordID, row)
}

func (s *ImportService) processFollowupRow(tx *gorm.DB, dryRun bool, batchID string, row map[string]string) error {
	parentKey := strings.TrimSpace(row["PARENT_KEY"])
	if parentKey == "" {
		return errors.New("PARENT_KEY is required")
	}
	key := strings.TrimSpace(row["KEY"])
	if key != "" {
		var existing int64
		if err := tx.Table("follow_up_actions").Where("source_key = ?", key).Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			return errors.New("duplicate follow-up KEY")
		}
	}
	if dryRun {
		return nil
	}
	return tx.Exec(`
		INSERT INTO follow_up_actions (
		  id, upload_batch_id, source_key, parent_key, linked_record_id, domain, challenge, intervention, responsibility, timelines
		) VALUES (
		  gen_random_uuid(), ?, NULLIF(?, ''), ?, (SELECT id FROM supervision_records WHERE source_key = ? LIMIT 1), ?, ?, ?, ?, ?
		)
	`, batchID, key, parentKey, parentKey, row["domain"], row["challenge"], row["intervention"], row["responsibility"], row["timelines"]).Error
}

func (s *ImportService) ensureRegion(tx *gorm.DB, name string) (string, error) {
	var id string
	err := tx.Raw(`SELECT id FROM regions WHERE name = ?`, name).Scan(&id).Error
	if err != nil {
		return "", err
	}
	if id != "" {
		return id, nil
	}
	if err := tx.Raw(`INSERT INTO regions (id, name) VALUES (gen_random_uuid(), ?) RETURNING id`, name).Scan(&id).Error; err != nil {
		return "", err
	}
	return id, nil
}

func (s *ImportService) ensureDistrict(tx *gorm.DB, regionID, name string) (string, error) {
	var id string
	err := tx.Raw(`SELECT id FROM districts WHERE region_id = ? AND name = ?`, regionID, name).Scan(&id).Error
	if err != nil {
		return "", err
	}
	if id != "" {
		return id, nil
	}
	if err := tx.Raw(`INSERT INTO districts (id, region_id, name) VALUES (gen_random_uuid(), ?, ?) RETURNING id`, regionID, name).Scan(&id).Error; err != nil {
		return "", err
	}
	return id, nil
}

func (s *ImportService) ensureFacility(tx *gorm.DB, regionID, districtID, name, level string) (string, error) {
	var id string
	err := tx.Raw(`SELECT id FROM facilities WHERE region_id = ? AND district_id = ? AND name = ?`, regionID, districtID, name).Scan(&id).Error
	if err != nil {
		return "", err
	}
	if id != "" {
		return id, nil
	}
	if err := tx.Raw(`INSERT INTO facilities (id, region_id, district_id, name, level) VALUES (gen_random_uuid(), ?, ?, ?, ?) RETURNING id`, regionID, districtID, name, level).Scan(&id).Error; err != nil {
		return "", err
	}
	return id, nil
}

func (s *ImportService) logImportError(tx *gorm.DB, batchID string, rowNum int, sourceColumn, code, detail string) error {
	return tx.Exec(`
		INSERT INTO import_errors (id, upload_batch_id, row_number, source_column, error_code, error_detail)
		VALUES (gen_random_uuid(), ?, ?, ?, ?, ?)
	`, batchID, rowNum, sourceColumn, code, detail).Error
}

func (s *ImportService) normalizeMainRow(tx *gorm.DB, recordID string, row map[string]string) error {
	baseCols := map[string]bool{
		"SubmissionDate":  true,
		"facility_name":   true,
		"level":           true,
		"region":          true,
		"district":        true,
		"period":          true,
		"meta-instanceID": true,
	}
	for col, raw := range row {
		if baseCols[col] {
			continue
		}
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		questionID, err := s.ensureQuestion(tx, col)
		if err != nil {
			return err
		}

		if isCommentColumn(col) {
			if err := tx.Exec(`
				INSERT INTO comments (id, record_id, question_id, text_value, created_at)
				VALUES (gen_random_uuid(), ?, ?, ?, NOW())
			`, recordID, questionID, raw).Error; err != nil {
				return err
			}
			continue
		}

		normalized, score, missing := normalizeScore(raw)
		if err := tx.Exec(`
			INSERT INTO responses_long (id, record_id, question_id, raw_value, normalized_value, numeric_score, is_missing)
			VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?)
			ON CONFLICT (record_id, question_id) DO UPDATE SET
				raw_value = EXCLUDED.raw_value,
				normalized_value = EXCLUDED.normalized_value,
				numeric_score = EXCLUDED.numeric_score,
				is_missing = EXCLUDED.is_missing
		`, recordID, questionID, raw, normalized, score, missing).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *ImportService) ensureQuestion(tx *gorm.DB, sourceColumn string) (string, error) {
	var id string
	if err := tx.Raw(`SELECT id FROM questions WHERE source_column = ?`, sourceColumn).Scan(&id).Error; err != nil {
		return "", err
	}
	if id != "" {
		return id, nil
	}
	qType := "text"
	scoreable := false
	if isCommentColumn(sourceColumn) {
		qType = "comment"
	} else if looksBinaryColumn(sourceColumn) {
		qType = "select_one yes_no"
		scoreable = true
	}
	if err := tx.Raw(`
		INSERT INTO questions (id, source_column, label, question_type, scoreable, option_list, display_order)
		VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, 0)
		RETURNING id
	`, sourceColumn, sourceColumn, qType, scoreable, nullableOptionList(scoreable)).Scan(&id).Error; err != nil {
		return "", err
	}
	return id, nil
}

func (s *ImportService) ListBatches(ctx context.Context, limit int) ([]UploadBatchDTO, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var rows []uploadBatchModel
	if err := s.repos.DB.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]UploadBatchDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, UploadBatchDTO{
			ID:            r.ID,
			BatchCode:     r.BatchCode,
			FileType:      r.FileType,
			FileName:      r.FileName,
			FileHash:      r.FileHash,
			Uploader:      r.Uploader,
			DryRun:        r.DryRun,
			Status:        r.Status,
			TotalRows:     r.TotalRows,
			ImportedRows:  r.ImportedRows,
			DuplicateRows: r.DuplicateRows,
			ErrorRows:     r.ErrorRows,
			CreatedAt:     r.CreatedAt,
			ProcessedAt:   r.ProcessedAt,
		})
	}
	return out, nil
}

func (s *ImportService) GetBatch(ctx context.Context, id string) (UploadBatchDTO, error) {
	var r uploadBatchModel
	if err := s.repos.DB.WithContext(ctx).Where("id = ?", id).First(&r).Error; err != nil {
		return UploadBatchDTO{}, err
	}
	return UploadBatchDTO{
		ID:            r.ID,
		BatchCode:     r.BatchCode,
		FileType:      r.FileType,
		FileName:      r.FileName,
		FileHash:      r.FileHash,
		Uploader:      r.Uploader,
		DryRun:        r.DryRun,
		Status:        r.Status,
		TotalRows:     r.TotalRows,
		ImportedRows:  r.ImportedRows,
		DuplicateRows: r.DuplicateRows,
		ErrorRows:     r.ErrorRows,
		CreatedAt:     r.CreatedAt,
		ProcessedAt:   r.ProcessedAt,
	}, nil
}

func (s *ImportService) GetBatchErrors(ctx context.Context, id string) ([]ImportErrorDTO, error) {
	var rows []importErrorModel
	if err := s.repos.DB.WithContext(ctx).Where("upload_batch_id = ?", id).Order("row_number ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]ImportErrorDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, ImportErrorDTO{
			ID:           r.ID,
			RowNumber:    r.RowNumber,
			SourceColumn: r.SourceColumn,
			ErrorCode:    r.ErrorCode,
			ErrorDetail:  r.ErrorDetail,
			CreatedAt:    r.CreatedAt,
		})
	}
	return out, nil
}

func toRowMap(headers, row []string) map[string]string {
	out := make(map[string]string, len(headers))
	for i, h := range headers {
		if i < len(row) {
			out[h] = strings.TrimSpace(row[i])
		} else {
			out[h] = ""
		}
	}
	return out
}

func toJSON(row map[string]string) string {
	parts := make([]string, 0, len(row))
	for k, v := range row {
		parts = append(parts, `"`+escapeJSON(k)+`":"`+escapeJSON(v)+`"`)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func escapeJSON(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`, "\t", `\t`)
	return r.Replace(s)
}

func validateHeaders(fileType string, headers []string) error {
	expected := mainHeaders
	if fileType == "followup" {
		expected = followHeaders
	}
	if len(headers) < len(expected) {
		return errors.New("invalid header count")
	}
	for i, h := range expected {
		if headers[i] != h {
			return errors.New("invalid header: " + headers[i])
		}
	}
	return nil
}

type SummaryDTO struct {
	OverallCompliance   float64 `json:"overall_compliance"`
	FacilitiesAssessed  int64   `json:"facilities_assessed"`
	UnresolvedFollowups int64   `json:"unresolved_followups"`
}

type DomainScoreDTO struct {
	Domain      string  `json:"domain"`
	Numerator   float64 `json:"numerator"`
	Denominator float64 `json:"denominator"`
	Compliance  float64 `json:"compliance"`
}

type RankingDTO struct {
	Name       string  `json:"name"`
	Compliance float64 `json:"compliance"`
}

type TrendDTO struct {
	Period     string  `json:"period"`
	Compliance float64 `json:"compliance"`
}

type QuestionPerformanceDTO struct {
	Question     string  `json:"question"`
	Compliance   float64 `json:"compliance"`
	Observations int64   `json:"observations"`
}

type GapDTO struct {
	Question     string  `json:"question"`
	Compliance   float64 `json:"compliance"`
	Observations int64   `json:"observations"`
}

type CommentDTO struct {
	ID        string    `json:"id"`
	Domain    string    `json:"domain"`
	Question  string    `json:"question"`
	TextValue string    `json:"text_value"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowupDTO struct {
	ID             string    `json:"id"`
	Domain         string    `json:"domain"`
	Challenge      string    `json:"challenge"`
	Intervention   string    `json:"intervention"`
	Responsibility string    `json:"responsibility"`
	Timelines      string    `json:"timelines"`
	Status         string    `json:"status"`
	ParentKey      string    `json:"parent_key"`
	CreatedAt      time.Time `json:"created_at"`
}

type CountByLabelDTO struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type AnalyticsFilter struct {
	Domain    string
	Subdomain string
	Period    string
	Region    string
	District  string
	Facility  string
	Level     string
}

func (s *AnalyticsService) Summary(ctx context.Context) (SummaryDTO, error) {
	return s.SummaryFiltered(ctx, AnalyticsFilter{})
}

func (s *AnalyticsService) SummaryFiltered(ctx context.Context, f AnalyticsFilter) (SummaryDTO, error) {
	out := SummaryDTO{}
	where, args := buildResponseFilterClause(f)
	if err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT COALESCE(AVG(numeric_score),0) AS overall_compliance
		FROM responses_long
		JOIN questions q ON q.id = responses_long.question_id
		JOIN supervision_records sr ON sr.id = responses_long.record_id
		LEFT JOIN domains d ON d.id = q.domain_id
		WHERE responses_long.is_missing = FALSE `+where, args...).Scan(&out).Error; err != nil {
		return SummaryDTO{}, err
	}
	recordWhere, recordArgs := buildRecordOnlyFilterClause(f)
	if err := s.repos.DB.WithContext(ctx).Raw(`SELECT COUNT(DISTINCT facility_id) FROM supervision_records sr WHERE 1=1 `+recordWhere, recordArgs...).Scan(&out.FacilitiesAssessed).Error; err != nil {
		return SummaryDTO{}, err
	}
	if err := s.repos.DB.WithContext(ctx).Raw(`SELECT COUNT(*) FROM follow_up_actions WHERE COALESCE(status,'open') <> 'closed'`).Scan(&out.UnresolvedFollowups).Error; err != nil {
		return SummaryDTO{}, err
	}
	return out, nil
}

func (s *AnalyticsService) Overview(ctx context.Context) (map[string]any, error) {
	sum, err := s.Summary(ctx)
	if err != nil {
		return nil, err
	}
	domains, err := s.DomainScores(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"summary": sum,
		"domains": domains,
	}, nil
}

func (s *AnalyticsService) DomainScores(ctx context.Context) ([]DomainScoreDTO, error) {
	return s.DomainScoresFiltered(ctx, AnalyticsFilter{})
}

func (s *AnalyticsService) DomainScoresFiltered(ctx context.Context, f AnalyticsFilter) ([]DomainScoreDTO, error) {
	var out []DomainScoreDTO
	where, args := buildResponseFilterClause(f)
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  COALESCE(d.name, 'Unmapped') AS domain,
		  COALESCE(SUM(rl.numeric_score),0) AS numerator,
		  COUNT(rl.id)::float AS denominator,
		  CASE WHEN COUNT(rl.id)=0 THEN 0 ELSE COALESCE(SUM(rl.numeric_score),0) / COUNT(rl.id)::float END AS compliance
		FROM responses_long rl
		JOIN questions q ON q.id = rl.question_id
		JOIN supervision_records sr ON sr.id = rl.record_id
		LEFT JOIN domains d ON d.id = q.domain_id
		WHERE rl.is_missing = FALSE `+where+`
		GROUP BY COALESCE(d.name, 'Unmapped')
		ORDER BY compliance DESC
	`, args...).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) FacilityRanking(ctx context.Context) ([]RankingDTO, error) {
	return s.FacilityRankingFiltered(ctx, AnalyticsFilter{})
}

func (s *AnalyticsService) FacilityRankingFiltered(ctx context.Context, f AnalyticsFilter) ([]RankingDTO, error) {
	var out []RankingDTO
	where, args := buildResponseFilterClause(f)
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  COALESCE(f.name, 'Unknown') AS name,
		  CASE WHEN COUNT(rl.id)=0 THEN 0 ELSE COALESCE(SUM(rl.numeric_score),0) / COUNT(rl.id)::float END AS compliance
		FROM supervision_records sr
		LEFT JOIN facilities f ON f.id = sr.facility_id
		LEFT JOIN responses_long rl ON rl.record_id = sr.id
		LEFT JOIN questions q ON q.id = rl.question_id
		LEFT JOIN domains d ON d.id = q.domain_id
		WHERE (rl.id IS NULL OR rl.is_missing = FALSE) `+where+`
		GROUP BY COALESCE(f.name, 'Unknown')
		ORDER BY compliance DESC
		LIMIT 20
	`, args...).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) DistrictRanking(ctx context.Context) ([]RankingDTO, error) {
	var out []RankingDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  COALESCE(d.name, 'Unknown') AS name,
		  CASE WHEN COUNT(rl.id)=0 THEN 0 ELSE COALESCE(SUM(rl.numeric_score),0) / COUNT(rl.id)::float END AS compliance
		FROM supervision_records sr
		LEFT JOIN districts d ON d.id = sr.district_id
		LEFT JOIN responses_long rl ON rl.record_id = sr.id AND rl.is_missing = FALSE
		GROUP BY COALESCE(d.name, 'Unknown')
		ORDER BY compliance DESC
		LIMIT 20
	`).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) RegionRanking(ctx context.Context) ([]RankingDTO, error) {
	var out []RankingDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  COALESCE(r.name, 'Unknown') AS name,
		  CASE WHEN COUNT(rl.id)=0 THEN 0 ELSE COALESCE(SUM(rl.numeric_score),0) / COUNT(rl.id)::float END AS compliance
		FROM supervision_records sr
		LEFT JOIN regions r ON r.id = sr.region_id
		LEFT JOIN responses_long rl ON rl.record_id = sr.id AND rl.is_missing = FALSE
		GROUP BY COALESCE(r.name, 'Unknown')
		ORDER BY compliance DESC
		LIMIT 20
	`).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) Trends(ctx context.Context) ([]TrendDTO, error) {
	return s.TrendsFiltered(ctx, AnalyticsFilter{})
}

func (s *AnalyticsService) TrendsFiltered(ctx context.Context, f AnalyticsFilter) ([]TrendDTO, error) {
	var out []TrendDTO
	where, args := buildResponseFilterClause(f)
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  COALESCE(sr.period, 'Unknown') AS period,
		  CASE WHEN COUNT(rl.id)=0 THEN 0 ELSE COALESCE(SUM(rl.numeric_score),0) / COUNT(rl.id)::float END AS compliance
		FROM supervision_records sr
		LEFT JOIN responses_long rl ON rl.record_id = sr.id
		LEFT JOIN questions q ON q.id = rl.question_id
		LEFT JOIN domains d ON d.id = q.domain_id
		WHERE (rl.id IS NULL OR rl.is_missing = FALSE) `+where+`
		GROUP BY COALESCE(sr.period, 'Unknown')
		ORDER BY period
	`, args...).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) QuestionPerformance(ctx context.Context) ([]QuestionPerformanceDTO, error) {
	return s.QuestionPerformanceFiltered(ctx, AnalyticsFilter{})
}

func (s *AnalyticsService) QuestionPerformanceFiltered(ctx context.Context, f AnalyticsFilter) ([]QuestionPerformanceDTO, error) {
	var out []QuestionPerformanceDTO
	where, args := buildResponseFilterClause(f)
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  q.source_column AS question,
		  CASE WHEN COUNT(rl.id)=0 THEN 0 ELSE COALESCE(SUM(rl.numeric_score),0) / COUNT(rl.id)::float END AS compliance,
		  COUNT(rl.id) AS observations
		FROM questions q
		LEFT JOIN responses_long rl ON rl.question_id = q.id
		LEFT JOIN supervision_records sr ON sr.id = rl.record_id
		LEFT JOIN domains d ON d.id = q.domain_id
		WHERE (rl.id IS NULL OR rl.is_missing = FALSE) `+where+`
		GROUP BY q.source_column
		ORDER BY compliance ASC, observations DESC
		LIMIT 100
	`, args...).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) Gaps(ctx context.Context) ([]GapDTO, error) {
	var out []GapDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  q.source_column AS question,
		  CASE WHEN COUNT(rl.id)=0 THEN 0 ELSE COALESCE(SUM(rl.numeric_score),0) / COUNT(rl.id)::float END AS compliance,
		  COUNT(rl.id) AS observations
		FROM questions q
		LEFT JOIN responses_long rl ON rl.question_id = q.id AND rl.is_missing = FALSE
		GROUP BY q.source_column
		HAVING COUNT(rl.id) > 0
		ORDER BY compliance ASC, observations DESC
		LIMIT 10
	`).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) Comments(ctx context.Context, query string, limit int) ([]CommentDTO, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	var out []CommentDTO
	if strings.TrimSpace(query) == "" {
		err := s.repos.DB.WithContext(ctx).Raw(`
			SELECT
			  c.id,
			  COALESCE(d.name, 'Unmapped') AS domain,
			  COALESCE(q.source_column, '') AS question,
			  c.text_value,
			  c.created_at
			FROM comments c
			LEFT JOIN questions q ON q.id = c.question_id
			LEFT JOIN domains d ON d.id = q.domain_id
			ORDER BY c.created_at DESC
			LIMIT ?
		`, limit).Scan(&out).Error
		return out, err
	}
	q := "%" + strings.TrimSpace(query) + "%"
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  c.id,
		  COALESCE(d.name, 'Unmapped') AS domain,
		  COALESCE(q.source_column, '') AS question,
		  c.text_value,
		  c.created_at
		FROM comments c
		LEFT JOIN questions q ON q.id = c.question_id
		LEFT JOIN domains d ON d.id = q.domain_id
		WHERE c.text_value ILIKE ?
		ORDER BY c.created_at DESC
		LIMIT ?
	`, q, limit).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) Followups(ctx context.Context, limit int) ([]FollowupDTO, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	var out []FollowupDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT
		  id,
		  COALESCE(domain,'Unknown') AS domain,
		  COALESCE(challenge,'') AS challenge,
		  COALESCE(intervention,'') AS intervention,
		  COALESCE(responsibility,'') AS responsibility,
		  COALESCE(timelines,'') AS timelines,
		  COALESCE(status,'open') AS status,
		  parent_key,
		  created_at
		FROM follow_up_actions
		ORDER BY created_at DESC
		LIMIT ?
	`, limit).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) FollowupStatus(ctx context.Context) ([]CountByLabelDTO, error) {
	var out []CountByLabelDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT COALESCE(status, 'open') AS label, COUNT(*) AS count
		FROM follow_up_actions
		GROUP BY COALESCE(status, 'open')
		ORDER BY count DESC
	`).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) FollowupByDomain(ctx context.Context) ([]CountByLabelDTO, error) {
	var out []CountByLabelDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT COALESCE(domain, 'Unknown') AS label, COUNT(*) AS count
		FROM follow_up_actions
		GROUP BY COALESCE(domain, 'Unknown')
		ORDER BY count DESC
	`).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) FollowupByResponsibility(ctx context.Context) ([]CountByLabelDTO, error) {
	var out []CountByLabelDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		SELECT COALESCE(responsibility, 'Unassigned') AS label, COUNT(*) AS count
		FROM follow_up_actions
		GROUP BY COALESCE(responsibility, 'Unassigned')
		ORDER BY count DESC
	`).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) FollowupTimelineBuckets(ctx context.Context) ([]CountByLabelDTO, error) {
	var out []CountByLabelDTO
	err := s.repos.DB.WithContext(ctx).Raw(`
		WITH bucketed AS (
		  SELECT
		    CASE
		      WHEN NULLIF(TRIM(timelines), '') IS NULL THEN 'unspecified'
		      WHEN timelines ~ '^[0-9]{4}-[0-9]{2}-[0-9]{2}$' AND timelines::date < CURRENT_DATE THEN 'overdue'
		      WHEN timelines ~ '^[0-9]{4}-[0-9]{2}-[0-9]{2}$' AND timelines::date <= CURRENT_DATE + INTERVAL '30 days' THEN '0-30 days'
		      WHEN timelines ~ '^[0-9]{4}-[0-9]{2}-[0-9]{2}$' AND timelines::date <= CURRENT_DATE + INTERVAL '90 days' THEN '31-90 days'
		      WHEN timelines ~ '^[0-9]{4}-[0-9]{2}-[0-9]{2}$' THEN '90+ days'
		      ELSE 'textual timeline'
		    END AS label
		  FROM follow_up_actions
		)
		SELECT label, COUNT(*) AS count
		FROM bucketed
		GROUP BY label
		ORDER BY count DESC
	`).Scan(&out).Error
	return out, err
}

func (s *AnalyticsService) DownloadCSV(ctx context.Context) (string, error) {
	domainScores, err := s.DomainScores(ctx)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	w := csv.NewWriter(&b)
	if err := w.Write([]string{"domain", "numerator", "denominator", "compliance"}); err != nil {
		return "", err
	}
	for _, d := range domainScores {
		row := []string{
			d.Domain,
			fmt.Sprintf("%.4f", d.Numerator),
			fmt.Sprintf("%.0f", d.Denominator),
			fmt.Sprintf("%.4f", d.Compliance),
		}
		if err := w.Write(row); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return b.String(), nil
}

func isCommentColumn(col string) bool {
	c := strings.ToLower(col)
	return strings.Contains(c, "comment") || strings.Contains(c, "recommendation")
}

func looksBinaryColumn(col string) bool {
	c := strings.ToLower(col)
	return strings.Contains(c, "_q") || strings.Contains(c, "question") || strings.HasPrefix(c, "d")
}

func normalizeScore(raw string) (normalized string, score float64, missing bool) {
	v := strings.ToLower(strings.TrimSpace(raw))
	switch v {
	case "yes", "1", "true":
		return "yes", 1, false
	case "no", "0", "false":
		return "no", 0, false
	case "":
		return "", 0, true
	default:
		return v, 0, false
	}
}

func nullableOptionList(scoreable bool) any {
	if scoreable {
		return "yes_no"
	}
	return nil
}

func buildResponseFilterClause(f AnalyticsFilter) (string, []any) {
	parts := make([]string, 0, 7)
	args := make([]any, 0, 8)
	if strings.TrimSpace(f.Period) != "" {
		parts = append(parts, " AND sr.period = ?")
		args = append(args, strings.TrimSpace(f.Period))
	}
	if strings.TrimSpace(f.Domain) != "" {
		parts = append(parts, " AND (LOWER(COALESCE(d.name,'')) LIKE LOWER(?) OR LOWER(COALESCE(q.source_column,'')) LIKE LOWER(?))")
		d := "%" + strings.TrimSpace(f.Domain) + "%"
		args = append(args, d, d)
	}
	if strings.TrimSpace(f.Subdomain) != "" {
		needle := "%" + strings.ReplaceAll(strings.TrimSpace(f.Subdomain), "-", "%") + "%"
		parts = append(parts, " AND LOWER(COALESCE(q.source_column,'')) LIKE LOWER(?)")
		args = append(args, needle)
	}
	if strings.TrimSpace(f.Region) != "" {
		parts = append(parts, " AND sr.region_id IN (SELECT id FROM regions WHERE LOWER(name) = LOWER(?))")
		args = append(args, strings.TrimSpace(f.Region))
	}
	if strings.TrimSpace(f.District) != "" {
		parts = append(parts, " AND sr.district_id IN (SELECT id FROM districts WHERE LOWER(name) = LOWER(?))")
		args = append(args, strings.TrimSpace(f.District))
	}
	if strings.TrimSpace(f.Facility) != "" {
		parts = append(parts, " AND sr.facility_id IN (SELECT id FROM facilities WHERE LOWER(name) = LOWER(?))")
		args = append(args, strings.TrimSpace(f.Facility))
	}
	if strings.TrimSpace(f.Level) != "" {
		parts = append(parts, " AND LOWER(COALESCE(sr.level,'')) = LOWER(?)")
		args = append(args, strings.TrimSpace(f.Level))
	}
	return strings.Join(parts, ""), args
}

func buildRecordOnlyFilterClause(f AnalyticsFilter) (string, []any) {
	parts := make([]string, 0, 5)
	args := make([]any, 0, 5)
	if strings.TrimSpace(f.Period) != "" {
		parts = append(parts, " AND sr.period = ?")
		args = append(args, strings.TrimSpace(f.Period))
	}
	if strings.TrimSpace(f.Region) != "" {
		parts = append(parts, " AND sr.region_id IN (SELECT id FROM regions WHERE LOWER(name) = LOWER(?))")
		args = append(args, strings.TrimSpace(f.Region))
	}
	if strings.TrimSpace(f.District) != "" {
		parts = append(parts, " AND sr.district_id IN (SELECT id FROM districts WHERE LOWER(name) = LOWER(?))")
		args = append(args, strings.TrimSpace(f.District))
	}
	if strings.TrimSpace(f.Facility) != "" {
		parts = append(parts, " AND sr.facility_id IN (SELECT id FROM facilities WHERE LOWER(name) = LOWER(?))")
		args = append(args, strings.TrimSpace(f.Facility))
	}
	if strings.TrimSpace(f.Level) != "" {
		parts = append(parts, " AND LOWER(COALESCE(sr.level,'')) = LOWER(?)")
		args = append(args, strings.TrimSpace(f.Level))
	}
	return strings.Join(parts, ""), args
}
