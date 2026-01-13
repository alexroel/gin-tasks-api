// @title Tasks API
// @version 1.0
// @description API para la gestión de tareas con autenticación de usuarios utilizando JWT y Gin Framework.
// @termsOfService http://swagger.io/terms/

// @contact.name Alex Roel
// @contact.url http://www.alexroel.com
// @contact.email roel@roelcode.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Tipo de token JWT con el prefijo 'Bearer '

package main

import (
	"log"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/handler"
	"github.com/alexroel/gin-tasks-api/internal/middleware"
	"github.com/alexroel/gin-tasks-api/internal/repository"
	"github.com/alexroel/gin-tasks-api/internal/service"
	"github.com/gin-gonic/gin"

	_ "github.com/alexroel/gin-tasks-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Cargar configuración
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Error en configuración: ", err)
	}

	// Conectar a la base de datos
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	defer config.CloseDB()

	// Ejecutar migraciones
	if err := config.RunMigrations(); err != nil {
		log.Fatal(err)
	}

	// Registrar repositorios
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()

	// Registrar servicios
	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Registrar Handlers
	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)

	// Iniciar el servidor
	router := gin.Default()

	// Ruta de documentación Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware de autenticación
	authMiddleware := middleware.AuthMiddleware(config.AppConfig.JWTSecret)

	// Rutas de autenticación
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/signup", authHandler.SignUpHandler)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.GET("/profile", authMiddleware, authHandler.Profile)
		authRoutes.PUT("/profile", authMiddleware, authHandler.UpdateProfile)
		authRoutes.DELETE("/profile", authMiddleware, authHandler.DeleteAccount)
	}

	// Rutas de tareas (protegidas)
	taskRoutes := router.Group("/api/tasks")
	taskRoutes.Use(authMiddleware)
	{
		taskRoutes.POST("", taskHandler.Create)
		taskRoutes.GET("", taskHandler.GetAll)
		taskRoutes.GET("/:id", taskHandler.GetByID)
		taskRoutes.PUT("/:id", taskHandler.Update)
		taskRoutes.DELETE("/:id", taskHandler.Delete)
		taskRoutes.PATCH("/:id/status", taskHandler.ToggleStatus)
	}

	router.Run(config.AppConfig.Port)
}
