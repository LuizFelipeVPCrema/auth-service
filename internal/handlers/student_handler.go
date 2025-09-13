package handlers

import (
	"net/http"

	"study-manager-service/internal/models"
	"study-manager-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StudentHandler representa o handler de estudantes
type StudentHandler struct {
	studentService *services.StudentService
}

// NewStudentHandler cria um novo handler de estudantes
func NewStudentHandler(studentService *services.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// CreateStudentHandler cria um novo estudante (rota pública)
func (h *StudentHandler) CreateStudentHandler(c *gin.Context) {
	var req models.StudentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "dados inválidos",
			"code":    "INVALID_DATA",
			"message": err.Error(),
		})
		return
	}

	// Para registro público, precisamos do user_id no body
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "user_id não fornecido",
			"code":    "MISSING_USER_ID",
			"message": "Header X-User-ID é obrigatório para registro de estudante",
		})
		return
	}

	student, err := h.studentService.CreateStudent(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao criar estudante",
			"code":    "CREATE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, student)
}

// GetStudentByIDHandler busca um estudante por ID
func (h *StudentHandler) GetStudentByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID inválido",
			"code":    "INVALID_ID",
			"message": "ID deve ser um UUID válido",
		})
		return
	}

	student, err := h.studentService.GetStudentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "estudante não encontrado",
			"code":    "NOT_FOUND",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

// GetStudentByUserIDHandler busca um estudante por User ID
func (h *StudentHandler) GetStudentByUserIDHandler(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "user_id não fornecido",
			"code":    "MISSING_USER_ID",
			"message": "user_id é obrigatório",
		})
		return
	}

	student, err := h.studentService.GetStudentByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "estudante não encontrado",
			"code":    "NOT_FOUND",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

// GetAllStudentsHandler lista todos os estudantes
func (h *StudentHandler) GetAllStudentsHandler(c *gin.Context) {
	students, err := h.studentService.GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao buscar estudantes",
			"code":    "FETCH_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": students,
		"count": len(students),
	})
}

// UpdateStudentHandler atualiza um estudante
func (h *StudentHandler) UpdateStudentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID inválido",
			"code":    "INVALID_ID",
			"message": "ID deve ser um UUID válido",
		})
		return
	}

	var req models.StudentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "dados inválidos",
			"code":    "INVALID_DATA",
			"message": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")
	student, err := h.studentService.UpdateStudent(id, &req, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "estudante não encontrado" {
			status = http.StatusNotFound
		} else if err.Error() == "acesso negado: estudante não pertence ao usuário" {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"error":   "erro ao atualizar estudante",
			"code":    "UPDATE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

// DeleteStudentHandler remove um estudante
func (h *StudentHandler) DeleteStudentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID inválido",
			"code":    "INVALID_ID",
			"message": "ID deve ser um UUID válido",
		})
		return
	}

	userID := c.GetString("user_id")
	err = h.studentService.DeleteStudent(id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "estudante não encontrado" {
			status = http.StatusNotFound
		} else if err.Error() == "acesso negado: estudante não pertence ao usuário" {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"error":   "erro ao remover estudante",
			"code":    "DELETE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
