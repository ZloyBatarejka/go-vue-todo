package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"goTodo/backend/models"
	"goTodo/backend/repository"
)

type TodoHandler struct {
	repo repository.TodoRepository
}

func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

// CreateTodo godoc
// @Summary Create todo
// @Tags todos
// @Accept json
// @Produce json
// @Param request body models.CreateTodoRequest true "Create todo request"
// @Success 201 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Value == "" {
		respondWithError(w, http.StatusBadRequest, "Field 'value' is required")
		return
	}

	todo := &models.Todo{
		Value: req.Value,
	}

	if err := h.repo.Create(todo); err != nil {
		log.Printf("Error creating todo: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	respondWithJSON(w, http.StatusCreated, todo)
}

// GetAllTodos godoc
// @Summary Get all todos
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} models.ErrorResponse
// @Router /todos [get]
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.GetAll()
	if err != nil {
		log.Printf("Error getting todos: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get todos")
		return
	}

	respondWithJSON(w, http.StatusOK, todos)
}

// GetTodo godoc
// @Summary Get todo by id
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := h.repo.GetByID(id)
	if err != nil {
		log.Printf("Error getting todo: %v", err)
		respondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

// DeleteTodo godoc
// @Summary Delete todo by id
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		log.Printf("Error deleting todo: %v", err)
		respondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.MessageResponse{Message: "Todo deleted successfully"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.ErrorResponse{Error: message})
}
