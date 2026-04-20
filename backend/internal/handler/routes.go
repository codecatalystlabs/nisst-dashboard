package handler

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"nisst/internal/service"
)

func Register(app *fiber.App, s *service.Registry) {
	api := app.Group("/api/v1")
	u := NewUploadHandler(s.Import)
	m := NewMetadataHandler(s.Metadata)
	r := NewRecordHandler()
	a := NewAnalyticsHandler(s.Analytics)
	d := NewDashboardHandler(s.Analytics)

	api.Post("/uploads/main", u.UploadMain)
	api.Post("/uploads/followups", u.UploadFollowups)
	api.Get("/uploads", u.List)
	api.Get("/uploads/:id", u.GetByID)
	api.Get("/uploads/:id/errors", u.Errors)
	api.Post("/uploads/:id/reprocess", u.Reprocess)

	api.Get("/metadata/domains", m.Domains)
	api.Get("/metadata/domains/:id/questions", m.DomainQuestions)
	api.Get("/metadata/questions", m.Questions)
	api.Get("/metadata/questions/:id", m.QuestionByID)
	api.Get("/metadata/options", m.Options)
	api.Get("/metadata/facilities", m.Facilities)
	api.Get("/metadata/regions", m.Regions)
	api.Get("/metadata/districts", m.Districts)
	api.Get("/metadata/periods", m.Periods)
	api.Get("/metadata/levels", m.Levels)

	api.Get("/records", r.List)
	api.Get("/records/:id", r.Get)
	api.Get("/records/:id/followups", r.Followups)
	api.Get("/records/export", r.Export)

	api.Get("/analytics/summary", a.Summary)
	api.Get("/analytics/overview", a.Overview)
	api.Get("/analytics/domain-scores", a.DomainScores)
	api.Get("/analytics/subdomain-scores", a.SubdomainScores)
	api.Get("/analytics/question-performance", a.QuestionPerformance)
	api.Get("/analytics/facility-ranking", a.FacilityRanking)
	api.Get("/analytics/district-ranking", a.DistrictRanking)
	api.Get("/analytics/region-ranking", a.RegionRanking)
	api.Get("/analytics/trends", a.Trends)
	api.Get("/analytics/gaps", a.Gaps)
	api.Get("/analytics/comments", a.Comments)
	api.Get("/analytics/followups", a.Followups)
	api.Get("/analytics/followups/status", a.FollowupStatus)
	api.Get("/analytics/followups/by-domain", a.FollowupByDomain)
	api.Get("/analytics/followups/by-responsibility", a.FollowupByResponsibility)
	api.Get("/analytics/followups/timeline-buckets", a.FollowupTimelineBuckets)
	api.Get("/analytics/download", a.Download)

	api.Get("/dashboard/kpis", d.KPIs)
	api.Get("/dashboard/charts/domain-compliance", d.DomainCompliance)
	api.Get("/dashboard/charts/region-comparison", d.RegionComparison)
	api.Get("/dashboard/charts/district-comparison", d.DistrictComparison)
	api.Get("/dashboard/charts/facility-level-comparison", d.FacilityLevelComparison)
	api.Get("/dashboard/charts/question-heatmap", d.QuestionHeatmap)
	api.Get("/dashboard/charts/followup-burden", d.FollowupBurden)
	api.Get("/dashboard/tables/low-performing-facilities", d.LowPerformingFacilities)
	api.Get("/dashboard/tables/unresolved-followups", d.UnresolvedFollowups)
	api.Get("/dashboard/tables/recent-comments", d.RecentComments)

	app.Get("/swagger-spec.yaml", func(c *fiber.Ctx) error {
		specPath, err := resolveSwaggerSpecPath()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.Type("yaml").SendFile(specPath)
	})
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger-spec.yaml",
		DeepLinking: true,
	}))
}

func resolveSwaggerSpecPath() (string, error) {
	candidates := []string{
		"docs/swagger.yaml",
		"../docs/swagger.yaml",
		"../../docs/swagger.yaml",
		"../../../docs/swagger.yaml",
	}
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 6; i++ {
			candidates = append(candidates, filepath.Join(dir, "docs", "swagger.yaml"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	for _, p := range candidates {
		clean := filepath.Clean(p)
		if _, err := os.Stat(clean); err == nil {
			return clean, nil
		}
	}
	return "", os.ErrNotExist
}
