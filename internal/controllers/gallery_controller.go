/*
 * Project Name: controllers
 * File: gallery_controller.go
 * Created Date: Monday December 29th 2025
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/nova-cdn/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package controllers

import (
	"nova-cdn/internal/repositories"
	"nova-cdn/internal/service"
	"nova-cdn/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GalleryController struct {
	GalleryRepo    *repositories.GalleryRepository
	GenerateRepo   *repositories.GenerateRepository
	GalleryService service.GalleryService
}

func NewGalleryController(db *gorm.DB) *GalleryController {
	return &GalleryController{
		GalleryRepo:    repositories.NewGalleryRepository(db),
		GenerateRepo:   repositories.NewGenerateRepository(db),
		GalleryService: service.NewGalleryService(db),
	}
}

// Index godoc
// @Summary List galleries
// @Description Get a paginated list of galleries
// @Tags galleries
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param subject_id query string false "Subject ID"
// @Param subject_type query string false "Subject Type"
// @Param size query string false "Size (original, small, medium, large)"
// @Success 200 {object} utils.PaginatedResponse{data=[]GallerySwagger}
// @Failure 400 {object} utils.SimpleErrorResponse
// @Router /galleries [get]
// @Security BearerAuth
func (ctrl *GalleryController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	subject_id := c.Query("subject_id", "")
	subject_type := c.Query("subject_type", "")
	size := c.Query("size", "")

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}

	total, err := ctrl.GalleryRepo.Count(subject_id, subject_type, size)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count galleries")
	}

	galleries, err := ctrl.GalleryRepo.FindAllPaginated(page, perPage, subject_id, subject_type, size)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve galleries")
	}

	if len(galleries) < 1 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "No galleries found")
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
// @Param subject_id formData int false "Subject ID"
// @Param subject_type formData string false "Subject Type"
// @Param dir formData string false "Directory name (gallery, payment, item, etc.)" default(gallery)
// @Param description formData string false "Image description"
// @Param is_private formData boolean false "Set image as private" default(false)
// @Success 201 {object} utils.Response{data=[]GallerySwagger}
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/upload [post]
// @Security BearerAuth
func (ctrl *GalleryController) Upload(c *fiber.Ctx) error {
	return ctrl.GalleryService.Upload(c)
}

// Destroy godoc
// @Summary Delete a gallery item (Soft Delete)
// @Description Move a gallery item to trash
// @Tags galleries
// @Accept json
// @Produce json
// @Param id path int true "Gallery ID"
// @Success 200 {object} utils.SimpleResponse
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/{id} [delete]
// @Security BearerAuth
func (ctrl *GalleryController) Destroy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid gallery ID")
	}

	gallery, err := ctrl.GalleryRepo.FindByID(uint64(id), false)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gallery not found")
	}

	if err := ctrl.GalleryRepo.Delete(gallery); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete gallery")
	}

	return utils.SimpleSuccessResponse(c, "Gallery deleted successfully")
}

// ForceDelete godoc
// @Summary Permanently delete a gallery item
// @Description Permanently delete a gallery item and its physical files
// @Tags galleries
// @Accept json
// @Produce json
// @Param id path int true "Gallery ID"
// @Success 200 {object} utils.SimpleResponse
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/{id}/force [delete]
// @Security BearerAuth
func (ctrl *GalleryController) ForceDelete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid gallery ID")
	}

	gallery, err := ctrl.GalleryRepo.FindByID(uint64(id), true)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gallery not found")
	}

	if err := ctrl.GalleryRepo.ForceDelete(gallery); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete gallery")
	}

	utils.RemoveImageFiles(gallery.FilePath)

	return utils.SimpleSuccessResponse(c, "Gallery deleted successfully")
}

// Restore godoc
// @Summary Restore a gallery item
// @Description Restore a gallery item from the trash
// @Tags galleries
// @Accept json
// @Produce json
// @Param id path int true "Gallery ID"
// @Success 200 {object} utils.Response{data=GallerySwagger}
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/{id}/restore [post]
// @Security BearerAuth
func (ctrl *GalleryController) Restore(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid gallery ID")
	}

	gallery, err := ctrl.GalleryRepo.FindByID(uint64(id), true)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gallery not found")
	}

	if err := ctrl.GalleryRepo.Restore(gallery); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to restore gallery")
	}

	return utils.SuccessResponse(c, "Gallery restored successfully", gallery)
}

// Show godoc
// @Summary Show a gallery item
// @Description Show a gallery item
// @Tags galleries
// @Accept json
// @Produce json
// @Param id path int true "Gallery ID"
// @Success 200 {object} utils.Response{data=GallerySwagger}
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/{id} [get]
// @Security BearerAuth
func (ctrl *GalleryController) Show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid gallery ID")
	}

	gallery, err := ctrl.GalleryRepo.FindByID(uint64(id), false)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gallery not found")
	}

	return utils.SuccessResponse(c, "Gallery retrieved successfully", gallery)
}

// ShowByGroupCode godoc
// @Summary Show galleries by group code
// @Description Show galleries by group code
// @Tags galleries
// @Accept json
// @Produce json
// @Param group_code path string true "Group Code"
// @Param size query string false "Size (original, small, medium, large)"
// @Success 200 {object} utils.Response{data=GallerySwagger}
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/{group_code} [get]
// @Security BearerAuth
func (ctrl *GalleryController) ShowByGroupCode(c *fiber.Ctx) error {
	groupCode := c.Params("group_code")
	size := c.Query("size", "")

	galleries, err := ctrl.GalleryRepo.FindByGroupCode(groupCode, size)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gallery not found")
	}

	return utils.SuccessResponse(c, "Galleries retrieved successfully", galleries)
}

// DestroyByGroupCode godoc
// @Summary Soft delete galleries by group code
// @Description Soft delete all galleries with the specified group code
// @Tags galleries
// @Accept json
// @Produce json
// @Param group_code path string true "Group Code"
// @Param size query string false "Size (original, small, medium, large)"
// @Success 200 {object} utils.SimpleResponse
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/{group_code} [delete]
// @Security BearerAuth
func (ctrl *GalleryController) DestroyByGroupCode(c *fiber.Ctx) error {
	groupCode := c.Params("group_code")
	size := c.Query("size", "")

	if err := ctrl.GalleryRepo.DeleteByGroupCode(groupCode, size); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete galleries")
	}

	return utils.SimpleSuccessResponse(c, "Galleries deleted successfully")
}

// ForceDeleteByGroupCode godoc
// @Summary Permanently delete galleries by group code
// @Description Permanently delete all galleries with the specified group code
// @Tags galleries
// @Accept json
// @Produce json
// @Param group_code path string true "Group Code"
// @Param size query string false "Size (original, small, medium, large)"
// @Success 200 {object} utils.SimpleResponse
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /galleries/{group_code}/force [delete]
// @Security BearerAuth
func (ctrl *GalleryController) ForceDeleteByGroupCode(c *fiber.Ctx) error {
	groupCode := c.Params("group_code")
	size := c.Query("size", "")

	galleries, err := ctrl.GalleryRepo.ForceDeleteByGroupCode(groupCode, size)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete galleries")
	}

	for _, gallery := range galleries {
		utils.RemoveImageFiles(gallery.FilePath)
	}

	return utils.SimpleSuccessResponse(c, "Galleries deleted successfully")
}
