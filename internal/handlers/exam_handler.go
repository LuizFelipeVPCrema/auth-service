package handlers

import (
	"net/http"

	"study-manager-service/internal/models"
	"study-manager-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExamHandler representa o handler de provas/trabalhos
type ExamHandler struct {
	examService *services.ExamService
}

// NewExamHandler cria um novo handler de provas/trabalhos
func NewExamHandler(examService *services.ExamService) *ExamHandler {
	return &ExamHandler{
		examService: examService,
	}
}

// CreateExamHandler cria uma nova prova/trabalho
func (h *ExamHandler) CreateExamHandler(c *gin.Context) {
	var req models.ExamCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "dados inválidos",
			"code":    "INVALID_DATA",
			"message": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")
	exam, err := h.examService.CreateExam(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao criar prova/trabalho",
			"code":    "CREATE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, exam)
}

// GetExamsHandler lista provas/trabalhos do usuário
func (h *ExamHandler) GetExamsHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	exams, err := h.examService.GetExams(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao buscar provas/trabalhos",
			"code":    "FETCH_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  exams,
		"count": len(exams),
	})
}

// GetExamByIDHandler busca uma prova/trabalho por ID
func (h *ExamHandler) GetExamByIDHandler(c *gin.Context) {
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
	exam, err := h.examService.GetExamByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "prova/trabalho não encontrado",
			"code":    "NOT_FOUND",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, exam)
}

// GetExamDetailsHandler busca uma prova/trabalho com detalhes
func (h *ExamHandler) GetExamDetailsHandler(c *gin.Context) {
	// Implementar busca com detalhes completos
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "funcionalidade não implementada",
		"code":    "NOT_IMPLEMENTED",
		"message": "Busca de prova/trabalho com detalhes ainda não implementada",
	})
}

// UpdateExamHandler atualiza uma prova/trabalho
func (h *ExamHandler) UpdateExamHandler(c *gin.Context) {
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

	var req models.ExamUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "dados inválidos",
			"code":    "INVALID_DATA",
			"message": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")
	exam, err := h.examService.UpdateExam(id, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao atualizar prova/trabalho",
			"code":    "UPDATE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, exam)
}

// DeleteExamHandler remove uma prova/trabalho
func (h *ExamHandler) DeleteExamHandler(c *gin.Context) {
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
	err = h.examService.DeleteExam(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao remover prova/trabalho",
			"code":    "DELETE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetExamsBySubjectHandler lista provas/trabalhos de uma matéria
func (h *ExamHandler) GetExamsBySubjectHandler(c *gin.Context) {
	// Implementar busca por matéria
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "funcionalidade não implementada",
		"code":    "NOT_IMPLEMENTED",
		"message": "Busca de provas/trabalhos por matéria ainda não implementada",
	})
}
