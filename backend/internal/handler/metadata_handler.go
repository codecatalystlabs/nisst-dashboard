package handler

import "github.com/gofiber/fiber/v2"

type MetadataHandler struct{}
func NewMetadataHandler() *MetadataHandler { return &MetadataHandler{} }
func (h *MetadataHandler) Domains(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) DomainQuestions(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) Questions(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) QuestionByID(c *fiber.Ctx) error { return c.JSON(fiber.Map{"id": c.Params("id")}) }
func (h *MetadataHandler) Options(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) Facilities(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) Regions(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) Districts(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *MetadataHandler) Periods(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
