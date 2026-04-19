package handler

import "github.com/gofiber/fiber/v2"

type RecordHandler struct{}
func NewRecordHandler() *RecordHandler { return &RecordHandler{} }
func (h *RecordHandler) List(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *RecordHandler) Get(c *fiber.Ctx) error { return c.JSON(fiber.Map{"id": c.Params("id")}) }
func (h *RecordHandler) Followups(c *fiber.Ctx) error { return c.JSON(fiber.Map{"items": []string{}}) }
func (h *RecordHandler) Export(c *fiber.Ctx) error { return c.JSON(fiber.Map{"url": "/tmp/records.csv"}) }
