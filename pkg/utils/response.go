package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ValidationErrorResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type Meta struct {
	TotalRecords int64 `json:"total_records"`
	ItemsOnPage  int   `json:"items_on_page"`
	PerPage      int   `json:"per_page"`
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	HasMorePages bool  `json:"has_more_pages"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
	})
}

func ValidationError(c *fiber.Ctx, errors map[string][]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(ValidationErrorResponse{
		Success: false,
		Message: "Validation error",
		Errors:  errors,
	})
}

func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, page, limit int, total int64, itemsOnPage int) error {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Meta{
			TotalRecords: total,
			ItemsOnPage:  itemsOnPage,
			PerPage:      limit,
			CurrentPage:  page,
			TotalPages:   totalPages,
			HasMorePages: page < totalPages,
		},
	})
}
