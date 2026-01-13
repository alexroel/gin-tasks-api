package config

import (
	"fmt"
	"log"
	"time"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia global de la base de datos que se utilizará en la aplicación
var DB *gorm.DB

// ConnectDB establece la conexión con la base de datos PostgreSQL
func ConnectDB() error {
	dsn := AppConfig.URLDatabase

	// Configurar logger de GORM según el modo
	var gormLogger logger.Interface
	if AppConfig.GinMode == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	// Configurar connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error al obtener instancia SQL: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)           // Conexiones inactivas en el pool
	sqlDB.SetMaxOpenConns(100)          // Máximo de conexiones abiertas
	sqlDB.SetConnMaxLifetime(time.Hour) // Tiempo máximo de vida de una conexión

	DB = db
	log.Println("Conexión a la base de datos exitosa")
	return nil
}

// RunMigrations ejecuta las migraciones de la base de datos
func RunMigrations() error {
	err := DB.AutoMigrate(
		&domain.User{},
		&domain.Task{},
	)
	if err != nil {
		return fmt.Errorf("error al ejecutar las migraciones: %w", err)
	}
	log.Println("Migraciones ejecutadas correctamente")
	return nil
}

// CloseDB cierra la conexión con la base de datos
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error al obtener la instancia de la base de datos:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Error al cerrar la conexión:", err)
	}
}
