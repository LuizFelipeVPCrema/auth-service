package services

import (
	"fmt"

	"study-manager-service/internal/clients"
	"study-manager-service/internal/models"
	"study-manager-service/internal/repositories"

	"github.com/google/uuid"
)

// StudentService representa o serviço de estudantes
type StudentService struct {
	studentRepo *repositories.StudentRepository
	authClient  *clients.AuthClient
}

// NewStudentService cria um novo serviço de estudantes
func NewStudentService(studentRepo *repositories.StudentRepository, authClient *clients.AuthClient) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
		authClient:  authClient,
	}
}

// CreateStudent cria um novo estudante
func (s *StudentService) CreateStudent(req *models.StudentCreateRequest, userID string) (*models.StudentResponse, error) {
	// Verificar se já existe um estudante com este user_id
	exists, err := s.studentRepo.ExistsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência do estudante: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("já existe um estudante cadastrado para este usuário")
	}

	// Criar estudante
	student := &models.Student{
		ID:    uuid.New(),
		Name:  req.Name,
		Email: req.Email,
		UserID: userID,
	}

	if err := s.studentRepo.Create(student); err != nil {
		return nil, fmt.Errorf("erro ao criar estudante: %w", err)
	}

	response := student.ToResponse()
	return &response, nil
}

// GetStudentByUserID busca um estudante por User ID
func (s *StudentService) GetStudentByUserID(userID string) (*models.StudentResponse, error) {
	student, err := s.studentRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("estudante não encontrado: %w", err)
	}

	response := student.ToResponse()
	return &response, nil
}

// UpdateStudent atualiza um estudante
func (s *StudentService) UpdateStudent(id uuid.UUID, req *models.StudentUpdateRequest, userID string) (*models.StudentResponse, error) {
	// Buscar estudante
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("estudante não encontrado: %w", err)
	}

	// Verificar se o estudante pertence ao usuário
	if student.UserID != userID {
		return nil, fmt.Errorf("acesso negado: estudante não pertence ao usuário")
	}

	// Atualizar campos fornecidos
	if req.Name != "" {
		student.Name = req.Name
	}
	if req.Email != "" {
		student.Email = req.Email
	}

	if err := s.studentRepo.Update(student); err != nil {
		return nil, fmt.Errorf("erro ao atualizar estudante: %w", err)
	}

	response := student.ToResponse()
	return &response, nil
}

// DeleteStudent remove um estudante
func (s *StudentService) DeleteStudent(id uuid.UUID, userID string) error {
	// Buscar estudante
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("estudante não encontrado: %w", err)
	}

	// Verificar se o estudante pertence ao usuário
	if student.UserID != userID {
		return fmt.Errorf("acesso negado: estudante não pertence ao usuário")
	}

	if err := s.studentRepo.Delete(id); err != nil {
		return fmt.Errorf("erro ao remover estudante: %w", err)
	}

	return nil
}

// GetAllStudents busca todos os estudantes
func (s *StudentService) GetAllStudents() ([]models.StudentResponse, error) {
	students, err := s.studentRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estudantes: %w", err)
	}

	responses := make([]models.StudentResponse, len(students))
	for i, student := range students {
		responses[i] = student.ToResponse()
	}

	return responses, nil
}

// GetStudentByID busca um estudante por ID
func (s *StudentService) GetStudentByID(id uuid.UUID) (*models.StudentResponse, error) {
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("estudante não encontrado: %w", err)
	}

	response := student.ToResponse()
	return &response, nil
}
