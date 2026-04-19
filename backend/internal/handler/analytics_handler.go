package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"nisst/internal/service"
)

type AnalyticsHandler struct{ svc *service.AnalyticsService }

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}
func (h *AnalyticsHandler) Summary(c *fiber.Ctx) error {
	item, err := h.svc.SummaryFiltered(context.Background(), readAnalyticsFilter(c))
	if err != nil {
		return err
	}
	return c.JSON(item)
}
func (h *AnalyticsHandler) Overview(c *fiber.Ctx) error {
	item, err := h.svc.Overview(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(item)
}
func (h *AnalyticsHandler) DomainScores(c *fiber.Ctx) error {
	items, err := h.svc.DomainScoresFiltered(context.Background(), readAnalyticsFilter(c))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) SubdomainScores(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"items": []string{}})
}
func (h *AnalyticsHandler) QuestionPerformance(c *fiber.Ctx) error {
	items, err := h.svc.QuestionPerformanceFiltered(context.Background(), readAnalyticsFilter(c))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) FacilityRanking(c *fiber.Ctx) error {
	items, err := h.svc.FacilityRankingFiltered(context.Background(), readAnalyticsFilter(c))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) DistrictRanking(c *fiber.Ctx) error {
	items, err := h.svc.DistrictRanking(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) RegionRanking(c *fiber.Ctx) error {
	items, err := h.svc.RegionRanking(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) Trends(c *fiber.Ctx) error {
	items, err := h.svc.TrendsFiltered(context.Background(), readAnalyticsFilter(c))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) Gaps(c *fiber.Ctx) error {
	items, err := h.svc.Gaps(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) Comments(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	items, err := h.svc.Comments(context.Background(), c.Query("q"), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) Followups(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	items, err := h.svc.Followups(context.Background(), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) FollowupStatus(c *fiber.Ctx) error {
	items, err := h.svc.FollowupStatus(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) FollowupByDomain(c *fiber.Ctx) error {
	items, err := h.svc.FollowupByDomain(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) FollowupByResponsibility(c *fiber.Ctx) error {
	items, err := h.svc.FollowupByResponsibility(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) FollowupTimelineBuckets(c *fiber.Ctx) error {
	items, err := h.svc.FollowupTimelineBuckets(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
func (h *AnalyticsHandler) Download(c *fiber.Ctx) error {
	csvContent, err := h.svc.DownloadCSV(context.Background())
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, "text/csv")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="analytics_export.csv"`)
	return c.SendString(csvContent)
}

func readAnalyticsFilter(c *fiber.Ctx) service.AnalyticsFilter {
	return service.AnalyticsFilter{
		Domain:    c.Query("domain"),
		Subdomain: c.Query("subdomain"),
		Period:    c.Query("period"),
		Region:    c.Query("region"),
		District:  c.Query("district"),
		Facility:  c.Query("facility"),
		Level:     c.Query("level"),
	}
}
