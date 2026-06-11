package handler

import (
	"strings"

	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/storage"
	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	storage *storage.LocalStorage
}

func NewFileHandler(store *storage.LocalStorage) *FileHandler {
	return &FileHandler{storage: store}
}

func (h *FileHandler) Serve(c *fiber.Ctx) error {
	rel := strings.TrimPrefix(c.Params("*"), "/")
	if rel == "" {
		return response.NotFound(c, "file not found")
	}
	full, err := h.storage.FullPath(rel)
	if err != nil {
		return response.NotFound(c, "file not found")
	}
	return c.SendFile(full)
}
