package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kaanchinar/todo-app/middleware"
	"github.com/kaanchinar/todo-app/models"
	"github.com/kaanchinar/todo-app/store"
)

type TodoHandler struct {
	store *store.CompositeStore
}

func NewTodoHandler(s *store.CompositeStore) *TodoHandler {
	return &TodoHandler{store: s}
}

// List godoc
// @Summary List all todos
// @Description Get all todos for the authenticated user
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Todo
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /todos [get]

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	todos, err := h.store.GetTodosByUserID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch todos"})
		return
	}

	writeJSON(w, http.StatusOK, todos)
}

// Create godoc
// @Summary Create a new todo
// @Description Create a new todo item for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.TodoRequest true "Todo details"
// @Success 201 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /todos [post]

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req models.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "title is required"})
		return
	}

	todo := &models.Todo{
		Title:     req.Title,
		Completed: req.Completed,
		UserID:    userID,
	}
	if err := h.store.CreateTodo(r.Context(), todo); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create todo"})
		return
	}

	writeJSON(w, http.StatusCreated, todo)
}

// Get godoc
// @Summary Get a todo by ID
// @Description Get a single todo item by its ID
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /todos/{id} [get]

func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid todo id"})
		return
	}

	todo, exists := h.store.GetTodoByID(r.Context(), id)
	if !exists || todo.UserID != userID {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{Error: "todo not found"})
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

// Update godoc
// @Summary Update a todo
// @Description Update an existing todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Param request body models.TodoRequest true "Updated todo details"
// @Success 200 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /todos/{id} [put]

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid todo id"})
		return
	}

	todo, exists := h.store.GetTodoByID(r.Context(), id)
	if !exists || todo.UserID != userID {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{Error: "todo not found"})
		return
	}

	var req models.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	todo.Completed = req.Completed

	if err := h.store.UpdateTodo(r.Context(), todo); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update todo"})
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

// Delete godoc
// @Summary Delete a todo
// @Description Delete a todo item by its ID
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /todos/{id} [delete]

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid todo id"})
		return
	}

	todo, exists := h.store.GetTodoByID(r.Context(), id)
	if !exists || todo.UserID != userID {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{Error: "todo not found"})
		return
	}

	if err := h.store.DeleteTodo(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to delete todo"})
		return
	}

	writeJSON(w, http.StatusOK, models.MessageResponse{Message: "todo deleted"})
}
