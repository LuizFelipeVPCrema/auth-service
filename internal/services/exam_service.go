package services

import (
	"fmt"

	"study-manager-service/internal/clients"
	"study-manager-service/internal/models"
	"study-manager-service/internal/repositories"

	"github.com/google/uuid"
)

// ExamService representa o serviço de provas/trabalhos
type ExamService struct {
	examRepo  *repositories.ExamRepository
	authClient *clients.AuthClient
}

// NewExamService cria um novo serviço de provas/trabalhos
func NewExamService(examRepo *repositories.ExamRepository, authClient *clients.AuthClient) *ExamService {
	return &ExamService{
		examRepo:  examRepo,
		authClient: authClient,
	}
}

// CreateExam cria uma nova prova/trabalho
func (s *ExamService) CreateExam(req *models.ExamCreateRequest, userID string) (*models.ExamResponse, error) {
	exam := &models.Exam{
		ID:          uuid.New(),
		SubjectID:   req.SubjectID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Type:        req.Type,
		Status:      req.Status,
	}

	if exam.Type == "" {
		exam.Type = "exam"
	}
	if exam.Status == "" {
		exam.Status = "pending"
	}

	if err := s.examRepo.Create(exam); err != nil {
		return nil, fmt.Errorf("erro ao criar prova/trabalho: %w", err)
	}

	response := exam.ToResponse()
	return &response, nil
}

// GetExams busca provas/trabalhos do usuário
func (s *ExamService) GetExams(userID string) ([]models.ExamResponse, error) {
	exams, err := s.examRepo.GetByStudentID(uuid.MustParse("00000000-0000-0000-0000-000000000000"))
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar provas/trabalhos: %w", err)
	}

	responses := make([]models.ExamResponse, len(exams))
	for i, exam := range exams {
		responses[i] = exam.ToResponse()
	}

	return responses, nil
}

// GetExamByID busca uma prova/trabalho por ID
func (s *ExamService) GetExamByID(id uuid.UUID, userID string) (*models.ExamResponse, error) {
	exam, err := s.examRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("prova/trabalho não encontrado: %w", err)
	}

	response := exam.ToResponse()
	return &response, nil
}

// UpdateExam atualiza uma prova/trabalho
func (s *ExamService) UpdateExam(id uuid.UUID, req *models.ExamUpdateRequest, userID string) (*models.ExamResponse, error) {
	exam, err := s.examRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("prova/trabalho não encontrado: %w", err)
	}

	if req.Title != "" {
		exam.Title = req.Title
	}
	if req.Description != "" {
		exam.Description = req.Description
	}
	if req.DueDate != nil {
		exam.DueDate = req.DueDate
	}
	if req.Type != "" {
		exam.Type = req.Type
	}
	if req.Status != "" {
		exam.Status = req.Status
	}

	if err := s.examRepo.Update(exam); err != nil {
		return nil, fmt.Errorf("erro ao atualizar prova/trabalho: %w", err)
	}

	response := exam.ToResponse()
	return &response, nil
}

// DeleteExam remove uma prova/trabalho
func (s *ExamService) DeleteExam(id uuid.UUID, userID string) error {
	_, err := s.examRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("prova/trabalho não encontrado: %w", err)
	}

	if err := s.examRepo.Delete(id); err != nil {
		return fmt.Errorf("erro ao remover prova/trabalho: %w", err)
	}

	return nil
}
