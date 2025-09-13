package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudyContent representa um conteúdo de estudo no sistema
type StudyContent struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	SubjectID   uuid.UUID      `json:"subject_id" gorm:"type:uuid;not null;index"`
	Subject     Subject        `json:"subject,omitempty" gorm:"foreignKey:SubjectID"`
	ExamID      *uuid.UUID     `json:"exam_id" gorm:"type:uuid;index"`
	Exam        *Exam          `json:"exam,omitempty" gorm:"foreignKey:ExamID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	IsCompleted bool           `json:"is_completed" gorm:"default:false"`
	Order       int            `json:"order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// StudyContentCreateRequest representa os dados para criar um conteúdo de estudo
type StudyContentCreateRequest struct {
	SubjectID   uuid.UUID `json:"subject_id" binding:"required"`
	ExamID      *uuid.UUID `json:"exam_id"`
	Title       string    `json:"title" binding:"required,min=2,max=200"`
	Description string    `json:"description" binding:"omitempty,max=1000"`
	Order       int       `json:"order" binding:"omitempty,min=0"`
}

// StudyContentUpdateRequest representa os dados para atualizar um conteúdo de estudo
type StudyContentUpdateRequest struct {
	Title       string `json:"title" binding:"omitempty,min=2,max=200"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Order       int    `json:"order" binding:"omitempty,min=0"`
}

// StudyContentCompleteRequest representa os dados para marcar como concluído
type StudyContentCompleteRequest struct {
	IsCompleted bool `json:"is_completed"`
}

// StudyContentReorderRequest representa os dados para reordenar conteúdos
type StudyContentReorderRequest struct {
	ContentIDs []uuid.UUID `json:"content_ids" binding:"required,min=1"`
}

// StudyContentResponse representa a resposta de um conteúdo de estudo
type StudyContentResponse struct {
	ID          uuid.UUID  `json:"id"`
	SubjectID   uuid.UUID  `json:"subject_id"`
	ExamID      *uuid.UUID `json:"exam_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	Order       int        `json:"order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToResponse converte StudyContent para StudyContentResponse
func (sc *StudyContent) ToResponse() StudyContentResponse {
	return StudyContentResponse{
		ID:          sc.ID,
		SubjectID:   sc.SubjectID,
		ExamID:      sc.ExamID,
		Title:       sc.Title,
		Description: sc.Description,
		IsCompleted: sc.IsCompleted,
		Order:       sc.Order,
		CreatedAt:   sc.CreatedAt,
		UpdatedAt:   sc.UpdatedAt,
	}
}
