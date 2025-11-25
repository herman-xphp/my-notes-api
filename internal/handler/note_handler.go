package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/herman-xphp/my-notes-api/internal/dto"
	"github.com/herman-xphp/my-notes-api/internal/middleware"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"github.com/herman-xphp/my-notes-api/internal/utils"
	"github.com/herman-xphp/my-notes-api/pkg/response"
)

// NoteHandler handles note endpoins
type NoteHandler struct {
	noteService service.NoteService
	helper      *utils.HandlerHelper
}

// NewNoteHandler creates a new note handler
func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
		helper:      utils.NewHandlerHelper(),
	}
}

// Create handles note creation
func (h *NoteHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req dto.CreateNoteRequest

	if !h.helper.ParseAndValidate(c, &req) {
		return nil
	}

	// Call service
	result, err := h.noteService.Create(c.Context(), userID, req)
	if err != nil {
		return err
	}
	return response.Created(c, "Note created successfully", result)
}

// GetByID handles getting a note by ID
func (h *NoteHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse note ID
	noteID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid note ID", err.Error())
	}

	// CAll service
	result, err := h.noteService.GetByID(c.Context(), uint(noteID), userID)
	if err != nil {
		return err
	}

	return response.Success(c, "Note retrieved successfully", result)
}

// GetAll handles getting all notes with pagination
func (h *NoteHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req dto.NoteQueryRequest

	if !h.helper.ParseQuery(c, &req) {
		return nil
	}

	// Call service
	result, err := h.noteService.GetAll(c.Context(), userID, req)
	if err != nil {
		return err
	}

	return response.Success(c, "Notes retrieved successfully", result)
}

// Update handles note update
func (h *NoteHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req dto.UpdateNoteRequest

	// Parse note ID
	noteID, ok := h.helper.ParseID(c, "id")
	if !ok {
		return nil
	}

	// Parse request body and validate request
	if !h.helper.ParseAndValidate(c, &req) {
		return nil
	}
	// Call service
	result, err := h.noteService.Update(c.Context(), uint(noteID), userID, req)
	if err != nil {
		return err
	}

	return response.Success(c, "Note updated successfully", result)
}

// Delete handles note deletion (soft delete)
func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse note ID
	noteID, ok := h.helper.ParseID(c, "id")
	if !ok {
		return nil
	}

	// Call service
	if err := h.noteService.Delete(c.Context(), uint(noteID), userID); err != nil {
		return err
	}

	return response.Success(c, "Note deleted successfully", nil)
}

// Restore handles note restoration
func (h *NoteHandler) Restore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse note ID
	noteID, ok := h.helper.ParseID(c, "id")
	if !ok {
		return nil
	}
	// Call service
	result, err := h.noteService.Restore(c.Context(), uint(noteID), userID)
	if err != nil {
		return err
	}

	return response.Success(c, "Note restored successfully", result)
}
