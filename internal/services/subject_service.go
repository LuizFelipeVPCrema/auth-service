package services

import (
	"fmt"

	"study-manager-service/internal/clients"
	"study-manager-service/internal/models"
	"study-manager-service/internal/repositories"

	"github.com/google/uuid"
)

// SubjectService representa o serviço de matérias
type SubjectService struct {
	subjectRepo *repositories.SubjectRepository
	authClient  *clients.AuthClient
}

// NewSubjectService cria um novo serviço de matérias
func NewSubjectService(subjectRepo *repositories.SubjectRepository, authClient *clients.AuthClient) *SubjectService {
	return &SubjectService{
		subjectRepo: subjectRepo,
		authClient:  authClient,
	}
}

// CreateSubject cria uma nova matéria
func (s *SubjectService) CreateSubject(req *models.SubjectCreateRequest, userID string) (*models.SubjectResponse, error) {
	// Buscar estudante por user_id
	// Aqui você precisaria de um repositório de estudantes
	// Por simplicidade, vou assumir que existe um método para buscar por user_id
	
	subject := &models.Subject{
		ID:          uuid.New(),
		StudentID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"), // Placeholder
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.subjectRepo.Create(subject); err != nil {
		return nil, fmt.Errorf("erro ao criar matéria: %w", err)
	}

	response := subject.ToResponse()
	return &response, nil
}

// GetSubjects busca matérias do usuário
func (s *SubjectService) GetSubjects(userID string) ([]models.SubjectResponse, error) {
	// Implementar busca por user_id
	subjects, err := s.subjectRepo.GetByStudentID(uuid.MustParse("00000000-0000-0000-0000-000000000000"))
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar matérias: %w", err)
	}

	responses := make([]models.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = subject.ToResponse()
	}

	return responses, nil
}

// GetSubjectByID busca uma matéria por ID
func (s *SubjectService) GetSubjectByID(id uuid.UUID, userID string) (*models.SubjectResponse, error) {
	subject, err := s.subjectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("matéria não encontrada: %w", err)
	}

	// Verificar propriedade (implementar verificação de student_id)
	
	response := subject.ToResponse()
	return &response, nil
}

// UpdateSubject atualiza uma matéria
func (s *SubjectService) UpdateSubject(id uuid.UUID, req *models.SubjectUpdateRequest, userID string) (*models.SubjectResponse, error) {
	subject, err := s.subjectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("matéria não encontrada: %w", err)
	}

	// Verificar propriedade

	if req.Name != "" {
		subject.Name = req.Name
	}
	if req.Description != "" {
		subject.Description = req.Description
	}

	if err := s.subjectRepo.Update(subject); err != nil {
		return nil, fmt.Errorf("erro ao atualizar matéria: %w", err)
	}

	response := subject.ToResponse()
	return &response, nil
}

// DeleteSubject remove uma matéria
func (s *SubjectService) DeleteSubject(id uuid.UUID, userID string) error {
	_, err := s.subjectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("matéria não encontrada: %w", err)
	}

	// Verificar propriedade

	if err := s.subjectRepo.Delete(id); err != nil {
		return fmt.Errorf("erro ao remover matéria: %w", err)
	}

	return nil
}
