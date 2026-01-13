package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config contiene todas las variables de configuración de la aplicación
type Config struct {
	// Aplicación
	GinMode string
	Port    string

	// Base de datos
	URLDatabase string

	// JWT
	JWTSecret   string
	JWTExpireIn time.Duration
}

// AppConfig es la instancia global de configuración
var AppConfig *Config

// LoadConfig carga las variables de entorno y valida la configuración
func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}

	// Parsear duración de JWT
	jwtExpireStr := getEnv("JWT_EXPIRE_IN", "24h")
	jwtExpire, err := time.ParseDuration(jwtExpireStr)
	if err != nil {
		return errors.New("JWT_EXPIRE_IN tiene un formato inválido: " + jwtExpireStr)
	}

	AppConfig = &Config{
		// Aplicación
		GinMode: getEnv("GIN_MODE", "debug"),
		Port:    getEnv("PORT", "8080"),

		// Base de datos
		URLDatabase: getEnv("URL_DATABASE", ""),

		// JWT
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTExpireIn: jwtExpire,
	}

	// Validar configuración crítica
	if err := validateConfig(); err != nil {
		return err
	}

	log.Println("Configuración cargada correctamente")
	return nil
}

// validateConfig valida que las variables críticas estén configuradas
func validateConfig() error {
	if AppConfig.JWTSecret == "" {
		return errors.New("JWT_SECRET es requerido y no puede estar vacío")
	}
	if len(AppConfig.JWTSecret) < 10 {
		return errors.New("JWT_SECRET debe tener al menos 10 caracteres")
	}
	if AppConfig.URLDatabase == "" {
		return errors.New("URL_DATABASE es requerido y no puede estar vacío")
	}
	return nil
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
