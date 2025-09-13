package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Exam representa uma prova ou trabalho no sistema
type Exam struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	SubjectID   uuid.UUID      `json:"subject_id" gorm:"type:uuid;not null;index"`
	Subject     Subject        `json:"subject,omitempty" gorm:"foreignKey:SubjectID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	DueDate     *time.Time     `json:"due_date"`
	Type        string         `json:"type" gorm:"default:'exam'"`
	Status      string         `json:"status" gorm:"default:'pending'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ExamCreateRequest representa os dados para criar uma prova/trabalho
type ExamCreateRequest struct {
	SubjectID   uuid.UUID  `json:"subject_id" binding:"required"`
	Title       string     `json:"title" binding:"required,min=2,max=200"`
	Description string     `json:"description" binding:"omitempty,max=1000"`
	DueDate     *time.Time `json:"due_date"`
	Type        string     `json:"type" binding:"omitempty,oneof=exam assignment project quiz"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
}

// ExamUpdateRequest representa os dados para atualizar uma prova/trabalho
type ExamUpdateRequest struct {
	Title       string     `json:"title" binding:"omitempty,min=2,max=200"`
	Description string     `json:"description" binding:"omitempty,max=1000"`
	DueDate     *time.Time `json:"due_date"`
	Type        string     `json:"type" binding:"omitempty,oneof=exam assignment project quiz"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
}

// ExamResponse representa a resposta de uma prova/trabalho
type ExamResponse struct {
	ID          uuid.UUID  `json:"id"`
	SubjectID   uuid.UUID  `json:"subject_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ExamDetailsResponse representa uma prova/trabalho com detalhes completos
type ExamDetailsResponse struct {
	ExamResponse
	StudyContents []StudyContentResponse `json:"study_contents,omitempty"`
	Attachments   []AttachmentResponse   `json:"attachments,omitempty"`
	References    []ExamReferenceResponse `json:"references,omitempty"`
}

// ToResponse converte Exam para ExamResponse
func (e *Exam) ToResponse() ExamResponse {
	return ExamResponse{
		ID:          e.ID,
		SubjectID:   e.SubjectID,
		Title:       e.Title,
		Description: e.Description,
		DueDate:     e.DueDate,
		Type:        e.Type,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// ToDetailsResponse converte Exam para ExamDetailsResponse
func (e *Exam) ToDetailsResponse(contents []StudyContent, attachments []Attachment, references []ExamReference) ExamDetailsResponse {
	contentResponses := make([]StudyContentResponse, len(contents))
	for i, content := range contents {
		contentResponses[i] = content.ToResponse()
	}

	attachmentResponses := make([]AttachmentResponse, len(attachments))
	for i, attachment := range attachments {
		attachmentResponses[i] = attachment.ToResponse()
	}

	referenceResponses := make([]ExamReferenceResponse, len(references))
	for i, reference := range references {
		referenceResponses[i] = reference.ToResponse()
	}

	return ExamDetailsResponse{
		ExamResponse:   e.ToResponse(),
		StudyContents:  contentResponses,
		Attachments:    attachmentResponses,
		References:     referenceResponses,
	}
}
