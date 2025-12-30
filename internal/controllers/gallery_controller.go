package controllers

import (
	"fmt"
	"nova-cdn/internal/models"
	"nova-cdn/internal/repositories"
	"nova-cdn/pkg/utils"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GalleryController struct {
	repo *repositories.GalleryRepository
}

func NewGalleryController(repo *repositories.GalleryRepository) *GalleryController {
	return &GalleryController{repo: repo}
}

// Index godoc
// @Summary List galleries
// @Description Get a paginated list of galleries
// @Tags galleries
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]GallerySwagger}
// @Failure 400 {object} utils.Response
// @Router /galleries [get]
// @Security BearerAuth
func (ctrl *GalleryController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}

	total, err := ctrl.repo.Count()

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count galleries")
	}

	galleries, err := ctrl.repo.FindAllPaginated(page, perPage)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve galleries")
	}

	return utils.PaginatedSuccessResponse(c, "Galleries retrieved successfully", galleries, page, perPage, total, len(galleries))
}

// Upload godoc
// @Summary Upload image to gallery
// @Description Upload a new image to the gallery with optimized versions
// @Tags galleries
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload"
// @Param description formData string false "Image description"
// @Param is_private formData boolean false "Set image as private" default(false)
// @Success 201 {object} utils.Response{data=[]GallerySwagger}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /galleries/upload [post]
// @Security BearerAuth
func (ctrl *GalleryController) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded")
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed")
	}

	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "File size exceeds 10MB limit")
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := "images/gallery/" + newFileName
	fullPath := "public/" + filePath
	outputDir := "public/images/gallery"

	if err := c.SaveFile(file, fullPath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save file")
	}

	description := c.FormValue("description", "")
	isPrivate := c.FormValue("is_private", "false") == "true"
	userID := c.Locals("user_id").(uint)

	original := &models.Gallery{
		UserID:       userID,
		FileName:     newFileName,
		FilePath:     filePath,
		FileSize:     uint32(file.Size),
		Description:  description,
		IsPrivate:    isPrivate,
		HasOptimized: true,
	}

	if err := ctrl.repo.Create(original); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save original image")
	}

	processedImages, err := utils.ProcessImage(fullPath, outputDir, newFileName)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process image: "+err.Error())
	}

	var galleries []*models.Gallery
	galleries = append(galleries, original)

	if len(processedImages) > 0 {
		subjectType := "App\\Models\\Gallery"
		var processedGalleries []*models.Gallery

		for _, img := range processedImages {
			processedGalleries = append(processedGalleries, &models.Gallery{
				UserID:       userID,
				SubjectID:    &original.ID,
				SubjectType:  &subjectType,
				FileName:     img.FileName,
				FilePath:     img.FilePath,
				FileSize:     img.FileSize,
				Description:  description,
				IsPrivate:    isPrivate,
				HasOptimized: false,
			})
		}

		if err := ctrl.repo.CreateMany(processedGalleries); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save processed images")
		}

		galleries = append(galleries, processedGalleries...)
	}

	return utils.CreatedResponse(c, "Image uploaded successfully", galleries)
}

// Destroy godoc
// @Summary Delete a gallery item (Soft Delete)
// @Description Move a gallery item to trash
// @Tags galleries
// @Accept json
// @Produce json
// @Param id path int true "Gallery ID"
// @Success 200 {object} utils.Response{data=GallerySwagger}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /galleries/{id} [delete]
// @Security BearerAuth
func (ctrl *GalleryController) Destroy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid gallery ID")
	}

	gallery, err := ctrl.repo.FindByID(uint64(id), false)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gallery not found")
	}

	if err := ctrl.repo.Delete(gallery); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete gallery")
	}

	return utils.SuccessResponse(c, "Gallery deleted successfully", gallery)
}

// ForceDelete godoc
// @Summary Permanently delete a gallery item
// @Description Permanently delete a gallery item and its physical files
// @Tags galleries
// @Accept json
// @Produce json
// @Param id path int true "Gallery ID"
// @Success 200 {object} utils.Response{data=GallerySwagger}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /galleries/{id}/force [delete]
// @Security BearerAuth
func (ctrl *GalleryController) ForceDelete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid gallery ID")
	}

	gallery, err := ctrl.repo.FindByID(uint64(id), true)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gallery not found")
	}

	if err := ctrl.repo.ForceDelete(gallery); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete gallery")
	}

	utils.RemoveImageFiles(gallery.FilePath)

	return utils.SuccessResponse(c, "Gallery deleted successfully", gallery)
}
