package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"nisst/internal/service"
)

type DashboardHandler struct{ svc *service.AnalyticsService }

func NewDashboardHandler(svc *service.AnalyticsService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) KPIs(c *fiber.Ctx) error {
	kpis, err := h.svc.Summary(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(kpis)
}

func (h *DashboardHandler) DomainCompliance(c *fiber.Ctx) error {
	items, err := h.svc.DomainScores(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *DashboardHandler) RegionComparison(c *fiber.Ctx) error {
	items, err := h.svc.RegionRanking(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *DashboardHandler) DistrictComparison(c *fiber.Ctx) error {
	items, err := h.svc.DistrictRanking(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *DashboardHandler) FacilityLevelComparison(c *fiber.Ctx) error {
	items, err := h.svc.FacilityRanking(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *DashboardHandler) QuestionHeatmap(c *fiber.Ctx) error {
	items, err := h.svc.QuestionPerformance(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *DashboardHandler) FollowupBurden(c *fiber.Ctx) error {
	items, err := h.svc.FollowupByDomain(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *DashboardHandler) LowPerformingFacilities(c *fiber.Ctx) error {
	items, err := h.svc.FacilityRanking(context.Background())
	if err != nil {
		return err
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit <= 0 {
		limit = 10
	}
	if len(items) > limit {
		items = items[len(items)-limit:]
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *DashboardHandler) UnresolvedFollowups(c *fiber.Ctx) error {
	items, err := h.svc.Followups(context.Background(), 100)
	if err != nil {
		return err
	}
	out := make([]service.FollowupDTO, 0, len(items))
	for _, item := range items {
		if item.Status != "closed" {
			out = append(out, item)
		}
	}
	return c.JSON(fiber.Map{"items": out})
}

func (h *DashboardHandler) RecentComments(c *fiber.Ctx) error {
	items, err := h.svc.Comments(context.Background(), c.Query("q"), 20)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
