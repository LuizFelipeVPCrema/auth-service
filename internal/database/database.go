package database

import (
	"log"
	"os"
	"path/filepath"

	"study-manager-service/internal/config"
	"study-manager-service/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database representa a conexão com o banco de dados
type Database struct {
	DB *gorm.DB
}

// NewDatabase cria uma nova conexão com o banco de dados
func NewDatabase(cfg *config.Config) (*Database, error) {
	// Criar diretório se não existir
	dir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Configurar logger do GORM
	var gormLogger logger.Interface
	if cfg.Server.Mode == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Conectar ao banco SQLite
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// Configurar pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Printf("Conexão com banco de dados SQLite estabelecida: %s", cfg.Database.Path)

	return &Database{DB: db}, nil
}

// Close fecha a conexão com o banco de dados
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate executa as migrações do banco de dados
func (d *Database) Migrate() error {
	err := d.DB.AutoMigrate(
		&models.Student{},
		&models.Subject{},
		&models.Exam{},
		&models.StudyContent{},
		&models.Attachment{},
		&models.ExamReference{},
	)

	if err != nil {
		return err
	}

	log.Println("Migrações do banco de dados executadas com sucesso")
	return nil
}

// CreateIndexes cria índices para melhorar a performance
func (d *Database) CreateIndexes() error {
	// Índices para Student
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_students_user_id ON students(user_id)").Error; err != nil {
		return err
	}

	// Índices para Subject
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_subjects_student_id ON subjects(student_id)").Error; err != nil {
		return err
	}

	// Índices para Exam
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_exams_subject_id ON exams(subject_id)").Error; err != nil {
		return err
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_exams_due_date ON exams(due_date)").Error; err != nil {
		return err
	}

	// Índices para StudyContent
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_study_contents_subject_id ON study_contents(subject_id)").Error; err != nil {
		return err
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_study_contents_exam_id ON study_contents(exam_id)").Error; err != nil {
		return err
	}

	// Índices para Attachment
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_attachments_exam_id ON attachments(exam_id)").Error; err != nil {
		return err
	}

	// Índices para ExamReference
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_exam_references_exam_id ON exam_references(exam_id)").Error; err != nil {
		return err
	}

	log.Println("Índices do banco de dados criados com sucesso")
	return nil
}
