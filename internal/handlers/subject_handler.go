package handlers

import (
	"net/http"

	"study-manager-service/internal/models"
	"study-manager-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubjectHandler representa o handler de matérias
type SubjectHandler struct {
	subjectService *services.SubjectService
}

// NewSubjectHandler cria um novo handler de matérias
func NewSubjectHandler(subjectService *services.SubjectService) *SubjectHandler {
	return &SubjectHandler{
		subjectService: subjectService,
	}
}

// CreateSubjectHandler cria uma nova matéria
func (h *SubjectHandler) CreateSubjectHandler(c *gin.Context) {
	var req models.SubjectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "dados inválidos",
			"code":    "INVALID_DATA",
			"message": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")
	subject, err := h.subjectService.CreateSubject(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao criar matéria",
			"code":    "CREATE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, subject)
}

// GetSubjectsHandler lista matérias do usuário
func (h *SubjectHandler) GetSubjectsHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	subjects, err := h.subjectService.GetSubjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao buscar matérias",
			"code":    "FETCH_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  subjects,
		"count": len(subjects),
	})
}

// GetSubjectByIDHandler busca uma matéria por ID
func (h *SubjectHandler) GetSubjectByIDHandler(c *gin.Context) {
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
	subject, err := h.subjectService.GetSubjectByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "matéria não encontrada",
			"code":    "NOT_FOUND",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// GetSubjectWithExamsHandler busca uma matéria com suas provas/trabalhos
func (h *SubjectHandler) GetSubjectWithExamsHandler(c *gin.Context) {
	// Implementar busca com provas/trabalhos
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "funcionalidade não implementada",
		"code":    "NOT_IMPLEMENTED",
		"message": "Busca de matéria com provas/trabalhos ainda não implementada",
	})
}

// UpdateSubjectHandler atualiza uma matéria
func (h *SubjectHandler) UpdateSubjectHandler(c *gin.Context) {
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

	var req models.SubjectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "dados inválidos",
			"code":    "INVALID_DATA",
			"message": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")
	subject, err := h.subjectService.UpdateSubject(id, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao atualizar matéria",
			"code":    "UPDATE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// DeleteSubjectHandler remove uma matéria
func (h *SubjectHandler) DeleteSubjectHandler(c *gin.Context) {
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
	err = h.subjectService.DeleteSubject(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "erro ao remover matéria",
			"code":    "DELETE_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
