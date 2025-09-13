package repositories

import (
	"study-manager-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentRepository representa o repositório de estudantes
type StudentRepository struct {
	db *gorm.DB
}

// NewStudentRepository cria um novo repositório de estudantes
func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// Create cria um novo estudante
func (r *StudentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

// GetByID busca um estudante por ID
func (r *StudentRepository) GetByID(id uuid.UUID) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("id = ?", id).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetByUserID busca um estudante por User ID
func (r *StudentRepository) GetByUserID(userID string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("user_id = ?", userID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetAll busca todos os estudantes
func (r *StudentRepository) GetAll() ([]models.Student, error) {
	var students []models.Student
	err := r.db.Find(&students).Error
	return students, err
}

// Update atualiza um estudante
func (r *StudentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

// Delete remove um estudante (soft delete)
func (r *StudentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Student{}, id).Error
}

// ExistsByUserID verifica se existe um estudante com o User ID
func (r *StudentRepository) ExistsByUserID(userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Student{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}
