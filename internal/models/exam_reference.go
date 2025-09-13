package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExamReference representa uma referência bibliográfica no sistema
type ExamReference struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	ExamID      uuid.UUID      `json:"exam_id" gorm:"type:uuid;not null;index"`
	Exam        Exam           `json:"exam,omitempty" gorm:"foreignKey:ExamID"`
	Title       string         `json:"title" gorm:"not null"`
	Author      string         `json:"author"`
	Publisher   string         `json:"publisher"`
	Year        int            `json:"year"`
	ISBN        string         `json:"isbn"`
	URL         string         `json:"url"`
	Description string         `json:"description" gorm:"type:text"`
	Type        string         `json:"type" gorm:"default:'book'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ExamReferenceCreateRequest representa os dados para criar uma referência bibliográfica
type ExamReferenceCreateRequest struct {
	ExamID      uuid.UUID `json:"exam_id" binding:"required"`
	Title       string    `json:"title" binding:"required,min=2,max=200"`
	Author      string    `json:"author" binding:"omitempty,max=100"`
	Publisher   string    `json:"publisher" binding:"omitempty,max=100"`
	Year        int       `json:"year" binding:"omitempty,min=1000,max=3000"`
	ISBN        string    `json:"isbn" binding:"omitempty,max=20"`
	URL         string    `json:"url" binding:"omitempty,url"`
	Description string    `json:"description" binding:"omitempty,max=1000"`
	Type        string    `json:"type" binding:"omitempty,oneof=book article website video other"`
}

// ExamReferenceUpdateRequest representa os dados para atualizar uma referência bibliográfica
type ExamReferenceUpdateRequest struct {
	Title       string `json:"title" binding:"omitempty,min=2,max=200"`
	Author      string `json:"author" binding:"omitempty,max=100"`
	Publisher   string `json:"publisher" binding:"omitempty,max=100"`
	Year        int    `json:"year" binding:"omitempty,min=1000,max=3000"`
	ISBN        string `json:"isbn" binding:"omitempty,max=20"`
	URL         string `json:"url" binding:"omitempty,url"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Type        string `json:"type" binding:"omitempty,oneof=book article website video other"`
}

// ExamReferenceResponse representa a resposta de uma referência bibliográfica
type ExamReferenceResponse struct {
	ID          uuid.UUID `json:"id"`
	ExamID      uuid.UUID `json:"exam_id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	Year        int       `json:"year"`
	ISBN        string    `json:"isbn"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converte ExamReference para ExamReferenceResponse
func (er *ExamReference) ToResponse() ExamReferenceResponse {
	return ExamReferenceResponse{
		ID:          er.ID,
		ExamID:      er.ExamID,
		Title:       er.Title,
		Author:      er.Author,
		Publisher:   er.Publisher,
		Year:        er.Year,
		ISBN:        er.ISBN,
		URL:         er.URL,
		Description: er.Description,
		Type:        er.Type,
		CreatedAt:   er.CreatedAt,
		UpdatedAt:   er.UpdatedAt,
	}
}
