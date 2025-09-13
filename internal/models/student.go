package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Student representa um estudante no sistema
type Student struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	UserID    string         `json:"user_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// StudentCreateRequest representa os dados para criar um estudante
type StudentCreateRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// StudentUpdateRequest representa os dados para atualizar um estudante
type StudentUpdateRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}

// StudentResponse representa a resposta de um estudante
type StudentResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converte Student para StudentResponse
func (s *Student) ToResponse() StudentResponse {
	return StudentResponse{
		ID:        s.ID,
		Name:      s.Name,
		Email:     s.Email,
		UserID:    s.UserID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
