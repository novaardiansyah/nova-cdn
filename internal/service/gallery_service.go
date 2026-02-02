/*
 * Project Name: service
 * File: gallery_service.go
 * Created Date: Sunday January 25th 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/nova-cdn/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package service

import (
	"fmt"
	"mime/multipart"
	"nova-cdn/internal/dto"
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

	if err := s.validateFile(file); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	input, err := s.parseUploadInput(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	filePath, newFileName, err := s.saveFileToDisk(c, file, input.Dir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	fullPath := filepath.Join(dto.UploadDirBase, filePath)
	outputDir := filepath.Join(dto.UploadDirBase, "images", input.Dir)

	groupCode := utils.GetCode(s.GenerateRepo, "gallery_group", true)

	original := models.Gallery{
		UserID:       input.UserID,
		SubjectID:    input.SubjectID,
		SubjectType:  input.SubjectType,
		FileName:     newFileName,
		FilePath:     filePath,
		FileSize:     uint32(file.Size),
		Description:  input.Description,
		IsPrivate:    input.IsPrivate,
		Size:         "original",
		HasOptimized: true,
		GroupCode:    groupCode,
	}

	if err := s.GalleryRepo.Create(&original); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save original image metadata")
	}

	processedImages, err := utils.ProcessImage(fullPath, outputDir, newFileName)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process image: "+err.Error())
	}

	if len(processedImages) > 0 {
		processedGalleries := s.buildProcessedGalleries(input, processedImages, filePath, groupCode)

		if err := s.GalleryRepo.CreateMany(processedGalleries); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save processed images metadata")
		}
	}

	galleries, err := s.GalleryRepo.FindByGroupCode(groupCode, "")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get processed images metadata")
	}

	return utils.CreatedResponse(c, "Image uploaded successfully", galleries)
}

func (s *galleryService) validateFile(file *multipart.FileHeader) error {
	contentType := file.Header.Get("Content-Type")
	if !dto.AllowedMimeTypes[contentType] {
		return fmt.Errorf("invalid file type. Only JPEG, PNG, GIF, and WebP are allowed")
	}
	if file.Size > dto.MaxUploadSize {
		return fmt.Errorf("file size exceeds 10MB limit")
	}
	return nil
}

func (s *galleryService) parseUploadInput(c *fiber.Ctx) (*dto.UploadInput, error) {
	userID := c.Locals("user_id").(uint)

	input := &dto.UploadInput{
		Dir:         c.FormValue("dir", dto.DefaultImageDir),
		Description: c.FormValue("description", ""),
		IsPrivate:   c.FormValue("is_private", "false") == "true",
		UserID:      userID,
	}

	subjectIDStr := c.FormValue("subject_id", "")
	subjectTypeStr := c.FormValue("subject_type", "")

	if subjectIDStr != "" {
		sid, err := strconv.ParseUint(subjectIDStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid subject_id format")
		}
		usid := uint(sid)
		input.SubjectID = &usid

		if subjectTypeStr != "" {
			input.SubjectType = &subjectTypeStr
		} else {
			stype := dto.ModelPrefix + utils.ToCamelCase(input.Dir)
			input.SubjectType = &stype
		}
	}

	return input, nil
}

func (s *galleryService) saveFileToDisk(c *fiber.Ctx, file *multipart.FileHeader, dir string) (string, string, error) {
	ext := filepath.Ext(file.Filename)
	newUid, err := uuid.NewV7()

	if err != nil {
		return "", "", err
	}

	newFileName := fmt.Sprintf("%v%s", newUid.String(), ext)
	relativePath := fmt.Sprintf("images/%s/%s", dir, newFileName)
	fullPath := filepath.Join(dto.UploadDirBase, relativePath)
	outputDir := filepath.Join(dto.UploadDirBase, "images", dir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := c.SaveFile(file, fullPath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	return relativePath, newFileName, nil
}

func (s *galleryService) buildProcessedGalleries(input *dto.UploadInput, processedImages []utils.ProcessedImage, originalPath string, groupCode string) []*models.Gallery {
	var result []*models.Gallery
	dirPath := filepath.Dir(originalPath)

	for _, img := range processedImages {
		result = append(result, &models.Gallery{
			UserID:       input.UserID,
			SubjectID:    input.SubjectID,
			SubjectType:  input.SubjectType,
			FileName:     img.FileName,
			FilePath:     fmt.Sprintf("%s/%s", dirPath, img.FileName),
			FileSize:     img.FileSize,
			Description:  input.Description,
			IsPrivate:    input.IsPrivate,
			HasOptimized: false,
			Size:         img.Size,
			GroupCode:    groupCode,
		})
	}
	return result
}

// ! End Upload()
