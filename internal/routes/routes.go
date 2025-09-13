package routes

import (
	"net/http"

	"study-manager-service/internal/config"
	"study-manager-service/internal/handlers"
	"study-manager-service/internal/middleware"
	"study-manager-service/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(
	cfg *config.Config,
	authMiddleware *middleware.AuthMiddleware,
	rateLimiter *middleware.RateLimiter,
	validationMiddleware *middleware.ValidationMiddleware,
	auditLogger *middleware.AuditLogger,
	studentService *services.StudentService,
	subjectService *services.SubjectService,
	examService *services.ExamService,
) *gin.Engine {
	router := gin.Default()

	// Middlewares globais
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.ValidateOrigin(cfg))
	router.Use(rateLimiter.RateLimit())
	router.Use(validationMiddleware.ValidateInput())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Study Manager Service is running",
			"service": "study-manager-service",
		})
	})

	// Rotas públicas (apenas registro de estudante)
	public := router.Group("/api/v1")
	{
		public.POST("/students", studentService.CreateStudentHandler)
	}

	// Rotas protegidas
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.RequireClientID())
	protected.Use(authMiddleware.RequireAuth())
	protected.Use(auditLogger.AuditSensitiveOperations())
	{
			// Inicializar handlers
		studentHandler := handlers.NewStudentHandler(studentService)
		subjectHandler := handlers.NewSubjectHandler(subjectService)
		examHandler := handlers.NewExamHandler(examService)

		// Rotas de estudantes
		students := protected.Group("/students")
		{
			students.GET("", studentHandler.GetAllStudentsHandler)
			students.GET("/:id", studentHandler.GetStudentByIDHandler)
			students.GET("/user/:user_id", studentHandler.GetStudentByUserIDHandler)
			students.PUT("/:id", studentHandler.UpdateStudentHandler)
			students.DELETE("/:id", studentHandler.DeleteStudentHandler)
		}

		// Rotas de matérias
		subjects := protected.Group("/subjects")
		{
			subjects.POST("", subjectHandler.CreateSubjectHandler)
			subjects.GET("", subjectHandler.GetSubjectsHandler)
			subjects.GET("/:id", subjectHandler.GetSubjectByIDHandler)
			subjects.GET("/:id/exams", subjectHandler.GetSubjectWithExamsHandler)
			subjects.PUT("/:id", subjectHandler.UpdateSubjectHandler)
			subjects.DELETE("/:id", subjectHandler.DeleteSubjectHandler)
		}

		// Rotas de provas/trabalhos
		exams := protected.Group("/exams")
		{
			exams.POST("", examHandler.CreateExamHandler)
			exams.GET("", examHandler.GetExamsHandler)
			exams.GET("/:id", examHandler.GetExamByIDHandler)
			exams.GET("/:id/details", examHandler.GetExamDetailsHandler)
			exams.PUT("/:id", examHandler.UpdateExamHandler)
			exams.DELETE("/:id", examHandler.DeleteExamHandler)
		}

		// Rotas de provas/trabalhos por matéria
		protected.GET("/subjects/:subject_id/exams", examHandler.GetExamsBySubjectHandler)
	}

	return router
}
