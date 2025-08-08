package handlers

import (
	"net/http"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CreateTodoRequest true "Todo data"
// @Success 201 {object} domain.APIResponse{data=domain.Todo}
// @Failure 400 {object} domain.APIResponse{error=domain.APIError}
// @Failure 401 {object} domain.APIResponse{error=domain.APIError}
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req domain.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// Set user ID from context (from auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeUnauthorized,
				Message: "User not authenticated",
			},
		})
		return
	}

	req.UserID = userID.(uuid.UUID)

	todo, err := h.todoService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInternalError,
				Message: "Failed to create todo",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, domain.APIResponse{
		Success: true,
		Data:    todo,
	})
}

// GetTodos godoc
// @Summary Get todos
// @Description Get paginated list of todos
// @Tags todos
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} domain.APIResponse{data=domain.PaginatedResponse}
// @Failure 401 {object} domain.APIResponse{error=domain.APIError}
// @Router /todos [get]
func (h *TodoHandler) GetTodos(c *gin.Context) {
	var pagination domain.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// Get user's todos
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeUnauthorized,
				Message: "User not authenticated",
			},
		})
		return
	}

	resp, err := h.todoService.GetByUserID(c.Request.Context(), userID.(uuid.UUID), pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInternalError,
				Message: "Failed to get todos",
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// GetTodo godoc
// @Summary Get todo by ID
// @Description Get a single todo by ID
// @Tags todos
// @Security BearerAuth
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} domain.APIResponse{data=domain.Todo}
// @Failure 400 {object} domain.APIResponse{error=domain.APIError}
// @Failure 404 {object} domain.APIResponse{error=domain.APIError}
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: "Invalid todo ID",
			},
		})
		return
	}

	todo, err := h.todoService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeTodoNotFound,
				Message: "Todo not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data:    todo,
	})
}

// UpdateTodo godoc
// @Summary Update todo
// @Description Update a todo by ID
// @Tags todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param request body domain.UpdateTodoRequest true "Todo update data"
// @Success 200 {object} domain.APIResponse{data=domain.Todo}
// @Failure 400 {object} domain.APIResponse{error=domain.APIError}
// @Failure 404 {object} domain.APIResponse{error=domain.APIError}
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: "Invalid todo ID",
			},
		})
		return
	}

	var req domain.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		})
		return
	}

	todo, err := h.todoService.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeTodoNotFound,
				Message: "Todo not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data:    todo,
	})
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description Delete a todo by ID
// @Tags todos
// @Security BearerAuth
// @Produce json
// @Param id path string true "Todo ID"
// @Success 204
// @Failure 400 {object} domain.APIResponse{error=domain.APIError}
// @Failure 404 {object} domain.APIResponse{error=domain.APIError}
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: "Invalid todo ID",
			},
		})
		return
	}

	err = h.todoService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeTodoNotFound,
				Message: "Todo not found",
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllTodos godoc
// @Summary Get all todos (admin)
// @Description Get paginated list of all todos (admin endpoint)
// @Tags todos
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} domain.APIResponse{data=domain.PaginatedResponse}
// @Failure 401 {object} domain.APIResponse{error=domain.APIError}
// @Router /admin/todos [get]
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	var pagination domain.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		})
		return
	}

	resp, err := h.todoService.List(c.Request.Context(), pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInternalError,
				Message: "Failed to get todos",
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data:    resp,
	})
}
