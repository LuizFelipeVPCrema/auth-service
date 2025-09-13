package repositories

import (
	"study-manager-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExamRepository representa o repositório de provas/trabalhos
type ExamRepository struct {
	db *gorm.DB
}

// NewExamRepository cria um novo repositório de provas/trabalhos
func NewExamRepository(db *gorm.DB) *ExamRepository {
	return &ExamRepository{db: db}
}

// Create cria uma nova prova/trabalho
func (r *ExamRepository) Create(exam *models.Exam) error {
	return r.db.Create(exam).Error
}

// GetByID busca uma prova/trabalho por ID
func (r *ExamRepository) GetByID(id uuid.UUID) (*models.Exam, error) {
	var exam models.Exam
	err := r.db.Where("id = ?", id).First(&exam).Error
	if err != nil {
		return nil, err
	}
	return &exam, nil
}

// GetBySubjectID busca provas/trabalhos por Subject ID
func (r *ExamRepository) GetBySubjectID(subjectID uuid.UUID) ([]models.Exam, error) {
	var exams []models.Exam
	err := r.db.Where("subject_id = ?", subjectID).Find(&exams).Error
	return exams, err
}

// GetByStudentID busca provas/trabalhos por Student ID
func (r *ExamRepository) GetByStudentID(studentID uuid.UUID) ([]models.Exam, error) {
	var exams []models.Exam
	err := r.db.Joins("JOIN subjects ON exams.subject_id = subjects.id").
		Where("subjects.student_id = ?", studentID).
		Find(&exams).Error
	return exams, err
}

// Update atualiza uma prova/trabalho
func (r *ExamRepository) Update(exam *models.Exam) error {
	return r.db.Save(exam).Error
}

// Delete remove uma prova/trabalho (soft delete)
func (r *ExamRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Exam{}, id).Error
}
