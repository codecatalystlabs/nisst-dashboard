package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"nisst/internal/service"
)

type UploadHandler struct{ svc *service.ImportService }

func NewUploadHandler(s *service.ImportService) *UploadHandler { return &UploadHandler{svc: s} }

func (h *UploadHandler) UploadMain(c *fiber.Ctx) error      { return h.upload(c, "main") }
func (h *UploadHandler) UploadFollowups(c *fiber.Ctx) error { return h.upload(c, "followup") }

func (h *UploadHandler) upload(c *fiber.Ctx, fileType string) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	summary, err := h.svc.ProcessUpload(context.Background(), service.ImportInput{
		FileType: fileType,
		FileName: fh.Filename,
		Uploader: c.FormValue("uploader", "system"),
		DryRun:   c.QueryBool("dry_run", false),
		Reader:   f,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(summary)
}

func (h *UploadHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	items, err := h.svc.ListBatches(context.Background(), limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *UploadHandler) GetByID(c *fiber.Ctx) error {
	item, err := h.svc.GetBatch(context.Background(), c.Params("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "batch not found")
		}
		return err
	}
	return c.JSON(item)
}

func (h *UploadHandler) Errors(c *fiber.Ctx) error {
	items, err := h.svc.GetBatchErrors(context.Background(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *UploadHandler) Reprocess(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "queued", "id": c.Params("id")})
}
