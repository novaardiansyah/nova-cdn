package service

import (
	"fmt"
	"nova-cdn/internal/models"
	"nova-cdn/internal/repositories"
	"nova-cdn/pkg/utils"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GalleryService interface {
	Upload(c *fiber.Ctx) error
}

type galleryService struct {
	GalleryRepo  *repositories.GalleryRepository
	GenerateRepo *repositories.GenerateRepository
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryRepo:  repositories.NewGalleryRepository(db),
		GenerateRepo: repositories.NewGenerateRepository(db),
	}
}

func (s *galleryService) Upload(c *fiber.Ctx) error {
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

	dir := c.FormValue("dir", "gallery")

	ext := filepath.Ext(file.Filename)
	newUid, err := uuid.NewV7()

	newFileName := fmt.Sprintf("%v%s", newUid.String(), ext)
	filePath := fmt.Sprintf("images/%s/%s", dir, newFileName)
	fullPath := "public/" + filePath
	outputDir := fmt.Sprintf("public/images/%s", dir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create directory")
	}

	if err := c.SaveFile(file, fullPath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save file")
	}

	description := c.FormValue("description", "")
	isPrivate := c.FormValue("is_private", "false") == "true"
	userID := c.Locals("user_id").(uint)

	subjectIDStr := c.FormValue("subject_id", "")
	subjectTypeStr := c.FormValue("subject_type", "")

	var subjectID *uint
	var subjectType *string

	if subjectIDStr != "" {
		if sid, err := strconv.ParseUint(subjectIDStr, 10, 32); err == nil {
			usid := uint(sid)
			subjectID = &usid

			if subjectTypeStr != "" {
				subjectType = &subjectTypeStr
			} else {
				stype := "App\\Models\\" + utils.ToCamelCase(dir)
				subjectType = &stype
			}
		}
	}

	groupCode := utils.GetCode(s.GenerateRepo, "gallery_group", true)

	original := &models.Gallery{
		UserID:       userID,
		SubjectID:    subjectID,
		SubjectType:  subjectType,
		FileName:     newFileName,
		FilePath:     filePath,
		FileSize:     uint32(file.Size),
		Description:  description,
		IsPrivate:    isPrivate,
		Size:         "original",
		HasOptimized: true,
		GroupCode:    groupCode,
	}

	if err := s.GalleryRepo.Create(original); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save original image")
	}

	processedImages, err := utils.ProcessImage(fullPath, outputDir, newFileName)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process image: "+err.Error())
	}

	var galleries []*models.Gallery
	galleries = append(galleries, original)

	if len(processedImages) > 0 {
		var processedGalleries []*models.Gallery

		for _, img := range processedImages {
			processedGalleries = append(processedGalleries, &models.Gallery{
				UserID:       userID,
				SubjectID:    subjectID,
				SubjectType:  subjectType,
				FileName:     img.FileName,
				FilePath:     fmt.Sprintf("images/%s/%s", dir, img.FileName),
				FileSize:     img.FileSize,
				Description:  description,
				IsPrivate:    isPrivate,
				HasOptimized: false,
				Size:         img.Size,
				GroupCode:    groupCode,
			})
		}

		if err := s.GalleryRepo.CreateMany(processedGalleries); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save processed images")
		}

		galleries = append(galleries, processedGalleries...)
	}

	return utils.CreatedResponse(c, "Image uploaded successfully", galleries)
}
