package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subject representa uma matéria no sistema
type Subject struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	StudentID   uuid.UUID      `json:"student_id" gorm:"type:uuid;not null;index"`
	Student     Student        `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// SubjectCreateRequest representa os dados para criar uma matéria
type SubjectCreateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

// SubjectUpdateRequest representa os dados para atualizar uma matéria
type SubjectUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

// SubjectResponse representa a resposta de uma matéria
type SubjectResponse struct {
	ID          uuid.UUID `json:"id"`
	StudentID   uuid.UUID `json:"student_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SubjectWithExamsResponse representa uma matéria com suas provas/trabalhos
type SubjectWithExamsResponse struct {
	SubjectResponse
	Exams []ExamResponse `json:"exams,omitempty"`
}

// ToResponse converte Subject para SubjectResponse
func (s *Subject) ToResponse() SubjectResponse {
	return SubjectResponse{
		ID:          s.ID,
		StudentID:   s.StudentID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// ToResponseWithExams converte Subject para SubjectWithExamsResponse
func (s *Subject) ToResponseWithExams(exams []Exam) SubjectWithExamsResponse {
	examResponses := make([]ExamResponse, len(exams))
	for i, exam := range exams {
		examResponses[i] = exam.ToResponse()
	}

	return SubjectWithExamsResponse{
		SubjectResponse: s.ToResponse(),
		Exams:          examResponses,
	}
}
