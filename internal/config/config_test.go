package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// saveEnvVars guarda y limpia todas las variables de entorno relevantes
func saveEnvVars() map[string]string {
	vars := []string{
		"JWT_SECRET", "DB_NAME", "JWT_EXPIRE_IN", "DB_HOST",
		"DB_PORT", "DB_USER", "DB_PASSWORD", "PORT", "GIN_MODE",
	}
	saved := make(map[string]string)
	for _, v := range vars {
		saved[v] = os.Getenv(v)
		os.Unsetenv(v)
	}
	return saved
}

// restoreEnvVars restaura las variables de entorno
func restoreEnvVars(saved map[string]string) {
	for k, v := range saved {
		if v != "" {
			os.Setenv(k, v)
		} else {
			os.Unsetenv(k)
		}
	}
	AppConfig = nil
}

func TestLoadConfig_Success(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Configurar variables de entorno
	os.Setenv("JWT_SECRET", "test-secret-key-12345")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_EXPIRE_IN", "24h")

	err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, AppConfig)
	assert.Equal(t, "test-secret-key-12345", AppConfig.JWTSecret)
	assert.Equal(t, "test_db", AppConfig.DBName)
}

func TestLoadConfig_MissingJWTSecret(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Configurar sin JWT_SECRET
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_EXPIRE_IN", "24h")

	err := LoadConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET es requerido")
}

func TestLoadConfig_JWTSecretTooShort(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Configurar con JWT_SECRET muy corto
	os.Setenv("JWT_SECRET", "short")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_EXPIRE_IN", "24h")

	err := LoadConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET debe tener al menos 10 caracteres")
}

func TestLoadConfig_MissingDBName(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Configurar sin DB_NAME - usará el valor por defecto "tasks_db"
	os.Setenv("JWT_SECRET", "test-secret-key-12345")
	os.Setenv("JWT_EXPIRE_IN", "24h")

	err := LoadConfig()

	// No debe dar error porque DB_NAME tiene valor por defecto
	assert.NoError(t, err)
	assert.Equal(t, "tasks_db", AppConfig.DBName)
}

func TestLoadConfig_InvalidJWTExpireIn(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Configurar con formato inválido de duración
	os.Setenv("JWT_SECRET", "test-secret-key-12345")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_EXPIRE_IN", "invalid-duration")

	err := LoadConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_EXPIRE_IN tiene un formato inválido")
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Solo configurar los requeridos
	os.Setenv("JWT_SECRET", "test-secret-key-12345")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_EXPIRE_IN", "24h")

	err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "localhost", AppConfig.DBHost)
	assert.Equal(t, "3306", AppConfig.DBPort)
	assert.Equal(t, "8080", AppConfig.Port)
	assert.Equal(t, "debug", AppConfig.GinMode)
}

func TestLoadConfig_CustomDuration(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Configurar con duración personalizada
	os.Setenv("JWT_SECRET", "test-secret-key-12345")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_EXPIRE_IN", "48h")

	err := LoadConfig()

	assert.NoError(t, err)
	// 48 horas = 48 * 60 * 60 segundos * 1e9 nanosegundos
	expected := 48 * 60 * 60 * 1e9
	assert.Equal(t, int64(expected), int64(AppConfig.JWTExpireIn))
}

func TestGetEnv_ExistingVariable(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default")

	assert.Equal(t, "test_value", value)
}

func TestGetEnv_DefaultValue(t *testing.T) {
	os.Unsetenv("NON_EXISTENT_VAR")

	value := getEnv("NON_EXISTENT_VAR", "default_value")

	assert.Equal(t, "default_value", value)
}

func TestValidateConfig_EmptyJWTSecret(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	AppConfig = &Config{
		JWTSecret: "",
		DBName:    "test_db",
	}

	err := validateConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET es requerido")
}

func TestValidateConfig_ShortJWTSecret(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	AppConfig = &Config{
		JWTSecret: "short",
		DBName:    "test_db",
	}

	err := validateConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET debe tener al menos 10 caracteres")
}

func TestValidateConfig_EmptyDBName(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	AppConfig = &Config{
		JWTSecret: "valid-secret-key-12345",
		DBName:    "",
	}

	err := validateConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_NAME es requerido")
}

func TestValidateConfig_Success(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	AppConfig = &Config{
		JWTSecret: "valid-secret-key-12345",
		DBName:    "test_db",
	}

	err := validateConfig()

	assert.NoError(t, err)
}

// ========== Table-Driven Tests ==========

func TestLoadConfig_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, cfg *Config)
	}{
		{
			name: "configuración mínima válida",
			envVars: map[string]string{
				"JWT_SECRET":    "minimum-secret-key",
				"DB_NAME":       "test_db",
				"JWT_EXPIRE_IN": "24h",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "minimum-secret-key", cfg.JWTSecret)
				assert.Equal(t, "test_db", cfg.DBName)
			},
		},
		{
			name: "configuración completa",
			envVars: map[string]string{
				"JWT_SECRET":    "complete-secret-key-12345",
				"DB_NAME":       "complete_db",
				"DB_HOST":       "mysql.example.com",
				"DB_PORT":       "3307",
				"DB_USER":       "admin",
				"DB_PASSWORD":   "secret123",
				"PORT":          "9090",
				"GIN_MODE":      "release",
				"JWT_EXPIRE_IN": "48h",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "complete-secret-key-12345", cfg.JWTSecret)
				assert.Equal(t, "complete_db", cfg.DBName)
				assert.Equal(t, "mysql.example.com", cfg.DBHost)
				assert.Equal(t, "3307", cfg.DBPort)
				assert.Equal(t, "admin", cfg.DBUser)
				assert.Equal(t, "secret123", cfg.DBPassword)
				assert.Equal(t, "9090", cfg.Port)
				assert.Equal(t, "release", cfg.GinMode)
			},
		},
		{
			name: "JWT_SECRET faltante",
			envVars: map[string]string{
				"DB_NAME":       "test_db",
				"JWT_EXPIRE_IN": "24h",
			},
			wantErr:     true,
			errContains: "JWT_SECRET es requerido",
		},
		{
			name: "JWT_SECRET muy corto",
			envVars: map[string]string{
				"JWT_SECRET":    "short",
				"DB_NAME":       "test_db",
				"JWT_EXPIRE_IN": "24h",
			},
			wantErr:     true,
			errContains: "JWT_SECRET debe tener al menos 10 caracteres",
		},
		{
			name: "JWT_EXPIRE_IN con formato inválido",
			envVars: map[string]string{
				"JWT_SECRET":    "valid-secret-key-12345",
				"DB_NAME":       "test_db",
				"JWT_EXPIRE_IN": "invalid",
			},
			wantErr:     true,
			errContains: "JWT_EXPIRE_IN tiene un formato inválido",
		},
		{
			name: "duración en minutos",
			envVars: map[string]string{
				"JWT_SECRET":    "valid-secret-key-12345",
				"DB_NAME":       "test_db",
				"JWT_EXPIRE_IN": "30m",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				expected := int64(30 * 60 * 1e9) // 30 minutos en nanosegundos
				assert.Equal(t, expected, int64(cfg.JWTExpireIn))
			},
		},
		{
			name: "duración en segundos",
			envVars: map[string]string{
				"JWT_SECRET":    "valid-secret-key-12345",
				"DB_NAME":       "test_db",
				"JWT_EXPIRE_IN": "3600s",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				expected := int64(3600 * 1e9) // 3600 segundos en nanosegundos
				assert.Equal(t, expected, int64(cfg.JWTExpireIn))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saved := saveEnvVars()
			defer restoreEnvVars(saved)

			// Configurar variables de entorno
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			err := LoadConfig()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, AppConfig)
				}
			}
		})
	}
}

func TestValidateConfig_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "configuración válida",
			config: &Config{
				JWTSecret: "valid-secret-key-12345",
				DBName:    "test_db",
			},
			wantErr: false,
		},
		{
			name: "JWT_SECRET vacío",
			config: &Config{
				JWTSecret: "",
				DBName:    "test_db",
			},
			wantErr:     true,
			errContains: "JWT_SECRET es requerido",
		},
		{
			name: "JWT_SECRET con 9 caracteres",
			config: &Config{
				JWTSecret: "123456789",
				DBName:    "test_db",
			},
			wantErr:     true,
			errContains: "JWT_SECRET debe tener al menos 10 caracteres",
		},
		{
			name: "JWT_SECRET con exactamente 10 caracteres",
			config: &Config{
				JWTSecret: "1234567890",
				DBName:    "test_db",
			},
			wantErr: false,
		},
		{
			name: "DB_NAME vacío",
			config: &Config{
				JWTSecret: "valid-secret-key-12345",
				DBName:    "",
			},
			wantErr:     true,
			errContains: "DB_NAME es requerido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saved := saveEnvVars()
			defer restoreEnvVars(saved)

			AppConfig = tt.config

			err := validateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnv_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
		setEnv       bool
	}{
		{
			name:         "variable existente",
			envKey:       "TEST_EXISTING",
			envValue:     "existing_value",
			defaultValue: "default",
			expected:     "existing_value",
			setEnv:       true,
		},
		{
			name:         "variable no existente",
			envKey:       "TEST_NON_EXISTING",
			envValue:     "",
			defaultValue: "fallback_value",
			expected:     "fallback_value",
			setEnv:       false,
		},
		{
			name:         "variable vacía usa default",
			envKey:       "TEST_EMPTY",
			envValue:     "",
			defaultValue: "default_for_empty",
			expected:     "",
			setEnv:       true, // Se establece como vacía
		},
		{
			name:         "variable con espacios",
			envKey:       "TEST_SPACES",
			envValue:     "  value with spaces  ",
			defaultValue: "default",
			expected:     "  value with spaces  ",
			setEnv:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.envKey)
			if tt.setEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			result := getEnv(tt.envKey, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadConfig_AllDefaultValues(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	// Solo configurar los requeridos
	os.Setenv("JWT_SECRET", "test-secret-key-12345")
	os.Setenv("JWT_EXPIRE_IN", "24h")

	err := LoadConfig()

	assert.NoError(t, err)
	// Verificar todos los valores por defecto
	assert.Equal(t, "localhost", AppConfig.DBHost)
	assert.Equal(t, "3306", AppConfig.DBPort)
	assert.Equal(t, "root", AppConfig.DBUser)
	assert.Equal(t, "", AppConfig.DBPassword)
	assert.Equal(t, "tasks_db", AppConfig.DBName)
	assert.Equal(t, "8080", AppConfig.Port)
	assert.Equal(t, "debug", AppConfig.GinMode)
}

func TestLoadConfig_ProductionMode(t *testing.T) {
	saved := saveEnvVars()
	defer restoreEnvVars(saved)

	os.Setenv("JWT_SECRET", "production-secret-key-very-secure-12345")
	os.Setenv("DB_NAME", "production_db")
	os.Setenv("DB_HOST", "prod-db.example.com")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "prod_user")
	os.Setenv("DB_PASSWORD", "super_secure_password")
	os.Setenv("PORT", "80")
	os.Setenv("GIN_MODE", "release")
	os.Setenv("JWT_EXPIRE_IN", "1h")

	err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "release", AppConfig.GinMode)
	assert.Equal(t, "80", AppConfig.Port)
	assert.Equal(t, "prod-db.example.com", AppConfig.DBHost)
}
