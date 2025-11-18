package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Status     string         `json:"status"`
	Message    string         `json:"message"`
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

func Success(c *fiber.Ctx, message string, data any) error {
	return c.JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string, err any) error {
	return c.Status(statusCode).JSON(Response{
		Status:  "error",
		Message: message,
		Error:   err,
	})
}

func BadRequest(c *fiber.Ctx, message string, err any) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, nil)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, nil)
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

func InternalServerError(c *fiber.Ctx, message string, err any) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

func SuccessWithPagination(c *fiber.Ctx, message string, data any, meta PaginationMeta) error {
	return c.JSON(PaginatedResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		Pagination: meta,
	})
}
