package repositories

import (
	"study-manager-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubjectRepository representa o repositório de matérias
type SubjectRepository struct {
	db *gorm.DB
}

// NewSubjectRepository cria um novo repositório de matérias
func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

// Create cria uma nova matéria
func (r *SubjectRepository) Create(subject *models.Subject) error {
	return r.db.Create(subject).Error
}

// GetByID busca uma matéria por ID
func (r *SubjectRepository) GetByID(id uuid.UUID) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.Where("id = ?", id).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// GetByStudentID busca matérias por Student ID
func (r *SubjectRepository) GetByStudentID(studentID uuid.UUID) ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Where("student_id = ?", studentID).Find(&subjects).Error
	return subjects, err
}

// GetByIDWithExams busca uma matéria com suas provas/trabalhos
func (r *SubjectRepository) GetByIDWithExams(id uuid.UUID) (*models.Subject, []models.Exam, error) {
	var subject models.Subject
	err := r.db.Where("id = ?", id).First(&subject).Error
	if err != nil {
		return nil, nil, err
	}

	var exams []models.Exam
	err = r.db.Where("subject_id = ?", id).Find(&exams).Error
	if err != nil {
		return nil, nil, err
	}

	return &subject, exams, nil
}

// Update atualiza uma matéria
func (r *SubjectRepository) Update(subject *models.Subject) error {
	return r.db.Save(subject).Error
}

// Delete remove uma matéria (soft delete)
func (r *SubjectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Subject{}, id).Error
}

// ExistsByStudentID verifica se existe uma matéria com o Student ID
func (r *SubjectRepository) ExistsByStudentID(studentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Subject{}).Where("student_id = ?", studentID).Count(&count).Error
	return count > 0, err
}
