package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Attachment representa um anexo no sistema
type Attachment struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	ExamID       uuid.UUID      `json:"exam_id" gorm:"type:uuid;not null;index"`
	Exam         Exam           `json:"exam,omitempty" gorm:"foreignKey:ExamID"`
	FileName     string         `json:"file_name" gorm:"not null"`
	OriginalName string         `json:"original_name" gorm:"not null"`
	FilePath     string         `json:"file_path" gorm:"not null"`
	FileSize     int64          `json:"file_size"`
	MimeType     string         `json:"mime_type"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// AttachmentCreateRequest representa os dados para criar um anexo
type AttachmentCreateRequest struct {
	ExamID      uuid.UUID `json:"exam_id" binding:"required"`
	Description string    `json:"description" binding:"omitempty,max=500"`
}

// AttachmentUpdateRequest representa os dados para atualizar um anexo
type AttachmentUpdateRequest struct {
	Description string `json:"description" binding:"omitempty,max=500"`
}

// AttachmentResponse representa a resposta de um anexo
type AttachmentResponse struct {
	ID           uuid.UUID `json:"id"`
	ExamID       uuid.UUID `json:"exam_id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse converte Attachment para AttachmentResponse
func (a *Attachment) ToResponse() AttachmentResponse {
	return AttachmentResponse{
		ID:           a.ID,
		ExamID:       a.ExamID,
		FileName:     a.FileName,
		OriginalName: a.OriginalName,
		FileSize:     a.FileSize,
		MimeType:     a.MimeType,
		Description:  a.Description,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
