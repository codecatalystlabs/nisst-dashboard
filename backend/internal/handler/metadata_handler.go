package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"nisst/internal/service"
)

type MetadataHandler struct{ svc *service.MetadataService }

func NewMetadataHandler(svc *service.MetadataService) *MetadataHandler {
	return &MetadataHandler{svc: svc}
}
func (h *MetadataHandler) Domains(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) DomainQuestions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"items": []string{}})
}
func (h *MetadataHandler) Questions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"items": []string{}})
}
func (h *MetadataHandler) QuestionByID(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"id": c.Params("id")})
}
func (h *MetadataHandler) Options(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }

func (h *MetadataHandler) Facilities(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	items, err := h.svc.Facilities(context.Background(), c.Query("q"), c.Query("region"), c.Query("district"), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *MetadataHandler) Regions(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	items, err := h.svc.Regions(context.Background(), c.Query("q"), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *MetadataHandler) Districts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	items, err := h.svc.Districts(context.Background(), c.Query("q"), c.Query("region"), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *MetadataHandler) Periods(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	items, err := h.svc.Periods(context.Background(), c.Query("q"), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *MetadataHandler) Levels(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	items, err := h.svc.Levels(context.Background(), c.Query("q"), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
