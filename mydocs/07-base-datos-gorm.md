# Base de Datos y GORM 

1. [Introducción](#introducción)
2. [¿Qué es GORM?](#qué-es-gorm)
3. [Conexión a MySQL con GORM](#conexión-a-mysql-con-gorm)
4. [Configuración con variables de entorno](#configuración-con-variables-de-entorno)
5. [Modelos con GORM](#modelos-con-gorm)
6. [Migraciones automáticas](#migraciones-automáticas)
7. [Relaciones](#relaciones)
8. [CRUD completo con GORM](#crud-completo-con-gorm)
9. [CRUD de Tareas con GORM](#crud-de-tareas-con-gorm)
10. [Queries avanzadas](#queries-avanzadas)
11. [Soft Delete](#soft-delete)
12. [Configurar logger y connection pool en GORM](#configurar-logger-y-connection-pool-en-gorm)
13. [Conclusión](#conclusión)



---
## Introducción
Bienvenido a la sección de Base de Datos y GORM. En esta sección, aprenderás a utilizar GORM, un ORM (Object-Relational Mapping) para Go, para interactuar con bases de datos relacionales de manera eficiente y sencilla.


---
## ¿Qué es GORM?
GORM es una biblioteca ORM para el lenguaje de programación Go que facilita la interacción con bases de datos relacionales. Proporciona una abstracción sobre las operaciones SQL, permitiendo a los desarrolladores trabajar con estructuras de datos Go en lugar de escribir consultas SQL manualmente. GORM soporta múltiples bases de datos, incluyendo MySQL, PostgreSQL, SQLite y SQL Server, y ofrece características como migraciones automáticas, relaciones entre modelos, consultas avanzadas y manejo de transacciones.

---
## Conexión a MySQL con GORM
Para conectar tu aplicación Go a una base de datos MySQL utilizando GORM, sigue estos pasos:

1. **Instala GORM y el controlador MySQL**:
   Asegúrate de tener GORM y el controlador MySQL instalados en tu proyecto. Puedes hacerlo ejecutando:
   ```bash
   go get -u gorm.io/gorm
   go get -u gorm.io/driver/mysql
   ```

2. **Configura la conexión a la base de datos**:
   Utiliza el siguiente código para establecer la conexión a tu base de datos MySQL:
    ```go
    //internal/config/database.go
    package config

    import (
        "fmt"
        "log"
        "time"

        "github.com/alexroel/gin-tasks-api/internal/domain"
        "gorm.io/driver/mysql"
        "gorm.io/gorm"
        "gorm.io/gorm/logger"
    )

    // DB es la instancia global de la base de datos que se utilizará en la aplicación
    var DB *gorm.DB

    // ConnectDB establece la conexión con la base de datos MySQL
    func ConnectDB() error {
        dsn := "roelcode:123456s@tcp(127.0.0.1:3306)/tasks_db?charset=utf8mb4&parseTime=True&loc=Local"

        DB = db
        log.Println("Conexión a la base de datos exitosa")
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
    ```

    3. **Llama a la función de conexión**:
    Asegúrate de llamar a `ConnectDB()` al iniciar tu aplicación para establecer la conexión a la base de datos.

    ```go
    // cmd/main.go
    package main

    import (
        "log"

        "github.com/alexroel/gin-tasks-api/internal/config"
    )

    func main() {
        if err := config.ConnectDB(); err != nil {
            log.Fatalf("Error al conectar a la base de datos: %v", err)
        }
        defer config.CloseDB()

        // Resto de la lógica de la aplicación...
    }
    ```
Con estos pasos, tu aplicación Go estará conectada a una base de datos MySQL utilizando GORM, y estarás listo para comenzar a definir modelos y realizar operaciones CRUD.

---
## Configuración con variables de entorno
Para mejorar la seguridad y flexibilidad de tu aplicación, es recomendable utilizar variables de entorno para almacenar la configuración de la base de datos en lugar de codificarla directamente en el código fuente. Aquí te mostramos cómo hacerlo:

1. **Instala la biblioteca `godotenv`**:
   Esta biblioteca te permite cargar variables de entorno desde un archivo `.env`. Instálala ejecutando:
   ```bash
   go get github.com/joho/godotenv
   ```

2. **Crea un archivo `.env`**:
    En la raíz de tu proyecto, crea un archivo llamado `.env` y agrega las siguientes variables de entorno:
    ```env
    # Variables de la aplicación
    GIN_MODE=debug
    PORT=5050

    # Configuración de la base de datos
    DB_HOST=localhost
    DB_PORT=3306
    DB_USER=roelcode
    DB_PASSWORD=123456
    DB_NAME=tasks_db

    # Configuración de JWT
    JWT_SECRET=mi_secreto_jwt
    JWT_EXPIRE_IN=48h
    ```
3. **Carga las variables de entorno en tu aplicación**:
    Ahora, vamos a crear un archivo `config.go` para cargar las variables de entorno y utilizarlas en la configuración de la base de datos.
    ```go
    // internal/config/config.go
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
        DBHost     string
        DBPort     string
        DBUser     string
        DBPassword string
        DBName     string

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
            DBHost:     getEnv("DB_HOST", "localhost"),
            DBPort:     getEnv("DB_PORT", "3306"),
            DBUser:     getEnv("DB_USER", "root"),
            DBPassword: getEnv("DB_PASSWORD", ""),
            DBName:     getEnv("DB_NAME", "tasks_db"),

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
        if len(AppConfig.JWTSecret) < 16 {
            return errors.New("JWT_SECRET debe tener al menos 16 caracteres")
        }
        if AppConfig.DBName == "" {
            return errors.New("DB_NAME es requerido")
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
    ```
    4. **Actualiza la configuración de la base de datos para usar las variables de entorno**:
    Modifica la función `ConnectDB` para utilizar las variables de entorno cargadas desde el archivo `.env`.
    ```go
    // internal/config/database.go
    // ...
    func ConnectDB() error {
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            AppConfig.DBUser,
            AppConfig.DBPassword,
            AppConfig.DBHost,
            AppConfig.DBPort,
            AppConfig.DBName,
        )

        // Resto del código...
    }
    ```
Con estos cambios, tu aplicación ahora utilizará variables de entorno para la configuración de la base de datos, lo que mejora la seguridad y facilita la gestión de diferentes entornos (desarrollo, producción, etc.).

---
## Modelos con GORM
En GORM, los modelos son estructuras de Go que representan las tablas en la base de datos. Cada campo en la estructura corresponde a una columna en la tabla. GORM utiliza etiquetas (tags) para definir las propiedades de cada campo, como el tipo de dato, restricciones y relaciones. Aquí tienes un ejemplo de cómo definir un modelo para una entidad `User` y `Task`:
```go
// internal/domain/user.go
package domain

import "gorm.io/gorm"

// User representa la entidad de un usuario en el sistema
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FullName  string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla para User
func (User) TableName() string {
	return "users"
}

// UserCreate representa los datos necesarios para crear un usuario
type UserCreate struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserUpdate representa los datos necesarios para actualizar un usuario
type UserUpdate struct {
	FullName *string `json:"full_name,omitempty" binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=8"`
}

// UserLogin representa los datos necesarios para que un usuario inicie sesión
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse representa la respuesta de un usuario sin datos sensibles
type UserResponse struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// ToResponse convierte un User a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
```

```go
// internal/domain/task.go
package domain

import "gorm.io/gorm"

// Task representa la entidad de una tarea en el sistema
type Task struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Completed bool           `gorm:"default:false" json:"completed"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla para Task
func (Task) TableName() string {
	return "tasks"
}

// CreateTask representa los datos necesarios para crear una nueva tarea
type CreateTask struct {
	Title string `json:"title" binding:"required,min=1,max=200"`
}

// UpdateTask representa los datos necesarios para actualizar una tarea existente
type UpdateTask struct {
	Title     *string `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Completed *bool   `json:"completed,omitempty"`
}

// TaskResponse representa la respuesta de una tarea con datos del usuario
type TaskResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	UserID    uint   `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
// ToResponse convierte un Task a TaskResponse
func (t *Task) ToResponse() TaskResponse {
    return TaskResponse{
        ID:        t.ID,
        Title:     t.Title,
        Completed: t.Completed,
        UserID:    t.UserID,
        CreatedAt: t.CreatedAt,
        UpdatedAt: t.UpdatedAt,
    }
}
```

Ahora que tienes definidos los modelos `User` y `Task`, puedes utilizar GORM para migrar las tablas en la base de datos, realizar operaciones CRUD y manejar relaciones entre las entidades.


---
## Migraciones automáticas
GORM proporciona una función de migración automática que crea o actualiza las tablas en la base de datos basándose en los modelos definidos en tu código. Esto facilita la gestión del esquema de la base de datos sin necesidad de escribir scripts SQL manualmente. Para utilizar las migraciones automáticas, sigue estos pasos:

```go
// internal/config/database.go
// ...
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
```

Ahoara para ejecutar las migraciones, simplemente llama a la función `RunMigrations()` después de establecer la conexión a la base de datos:

```go
// cmd/main.go
package main
import (
    "log"

    "github.com/alexroel/gin-tasks-api/internal/config"
)

func main() {
    if err := config.ConnectDB(); err != nil {
        log.Fatalf("Error al conectar a la base de datos: %v", err)
    }
    defer config.CloseDB()

    // Ejecutar migraciones
    if err := config.RunMigrations(); err != nil {
        log.Fatalf("Error al ejecutar migraciones: %v", err)
    }

    // Resto de la lógica de la aplicación...
}
``` 

Ahora, cada vez que inicies tu aplicación, GORM verificará los modelos definidos y aplicará cualquier cambio necesario en la base de datos para mantener el esquema actualizado. Esto incluye la creación de nuevas tablas, la adición de columnas y la modificación de tipos de datos según sea necesario.

---
## Relaciones 
GORM soporta varios tipos de relaciones entre modelos, como uno a muchos, muchos a muchos y uno a uno. Aquí te mostramos cómo definir y utilizar estas relaciones en tus modelos.


**Relación Uno a Muchos (User y Task)**:
En este ejemplo, un usuario puede tener muchas tareas. La relación se define utilizando la clave foránea `UserID` en el modelo `Task`.
```go
// internal/domain/task.go
// Task representa la entidad de una tarea en el sistema
type Task struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Completed bool           `gorm:"default:false" json:"completed"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

- `User  User` `gorm:"foreignKey:UserID" json:"-"` : Este campo establece la relación con el modelo `User` utilizando la clave foránea `UserID`. La etiqueta `json:"-"` indica que este campo no se incluirá en las respuestas JSON para evitar exponer datos innecesarios.

Y en el modelo `User`, puedes definir una relación inversa si lo deseas:
```go
// User representa la entidad de un usuario en el sistema
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FullName  string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Tasks     []Task         `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```
- `Tasks     []Task`        `gorm:"foreignKey:UserID" json:"tasks,omitempty"` : Este campo define una relación uno a muchos desde el usuario hacia sus tareas. La etiqueta `json:"tasks,omitempty"` permite incluir las tareas en las respuestas JSON si es necesario.
---
## CRUD completo con GORM
Ahora que hemos cubierto los conceptos básicos de GORM, vamos a profundizar en cómo realizar operaciones CRUD (Crear, Leer, Actualizar, Eliminar) utilizando GORM en nuestras entidades `User`:

```go
// internal/repository/user_repository.go
package repository

import (
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"gorm.io/gorm"
)

// UserRepository define las operaciones de base de datos para usuarios
type UserRepository interface {
	Create(user *domain.User) error
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	ExistsByEmail(email string) (bool, error)
}

// userRepository implementa UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{db: config.DB}
}

// Create crea un nuevo usuario en la base de datos
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetAll obtiene todos los usuarios
func (r *userRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}

// GetByID obtiene un usuario por su ID
func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update actualiza un usuario existente
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete elimina un usuario por su ID (soft delete)
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
```

Con este repositorio, puedes realizar operaciones CRUD completas en la entidad `User` utilizando GORM. Puedes seguir un enfoque similar para implementar el repositorio de `Task` y manejar las operaciones CRUD para las tareas.

---
## CRUD de Tareas con GORM
Ahora, vamos usar la IA para implementar el CRUD completo para la entidad `Task` utilizando GORM.
A la IA le pedimos con el siguiente prompt:
"Implementa el CRUD completo para la entidad Task utilizando GORM en Go. Crea un repositorio llamado TaskRepository con métodos para Crear, Obtener todos, Obtener por ID, Actualizar y Eliminar tareas. Asegúrate de manejar errores y utilizar las convenciones de GORM."
```go
// internal/repository/task_repository.go
package repository

import (
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"gorm.io/gorm"
)

// TaskRepository define las operaciones de base de datos para tareas
type TaskRepository interface {
	Create(task *domain.Task) error
	GetAll() ([]domain.Task, error)
	GetByID(id uint) (*domain.Task, error)
	GetByUserID(userID uint) ([]domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uint) error
}

// taskRepository implementa TaskRepository
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository crea una nueva instancia de TaskRepository
func NewTaskRepository() TaskRepository {
	return &taskRepository{db: config.DB}
}

// Create crea una nueva tarea en la base de datos
func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

// GetAll obtiene todas las tareas
func (r *taskRepository) GetAll() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// GetByID obtiene una tarea por su ID
func (r *taskRepository) GetByID(id uint) (*domain.Task, error) {
	var task domain.Task
	err := r.db.First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}

// Update actualiza una tarea existente
func (r *taskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

// Delete elimina una tarea por su ID (soft delete)
func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Task{}, id).Error
}
```

Lo anterior completa el CRUD para la entidad `Task` utilizando GORM. Ahora puedes crear, leer, actualizar y eliminar tareas en la base de datos de manera eficiente.

---
## Queries avanzadas
GORM ofrece una variedad de métodos para realizar consultas avanzadas en la base de datos. Aquí te mostramos algunos ejemplos de cómo utilizar estas funcionalidades para realizar consultas más complejas.
```go
// internal/repository/user_repository.go
// ...
// GetByEmail obtiene un usuario por su email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
```

De la misma manera, puedes implementar consultas avanzadas en el repositorio de `Task` según tus necesidades específicas.

```go
// internal/repository/task_repository.go
// ...
// GetByUserID obtiene todas las tareas de un usuario
func (r *taskRepository) GetByUserID(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

// Actualiza el estado de una tarea
func (r *taskRepository) UpdateStatus(id uint, completed bool) error {
    return r.db.Model(&domain.Task{}).Where("id = ?", id).Update("completed", completed).Error
}
```

Con estos ejemplos, puedes ver cómo utilizar GORM para realizar consultas avanzadas y personalizadas en tus repositorios, lo que te permite manejar datos de manera más eficiente y flexible.

---
## Soft Delete
GORM soporta la funcionalidad de "soft delete", que permite marcar registros como eliminados sin borrarlos físicamente de la base de datos. Esto es útil para mantener un historial de datos y permitir la recuperación de registros eliminados si es necesario. Para implementar soft delete en tus modelos, debes agregar el campo `gorm.DeletedAt` a la estructura del modelo. Aquí tienes un ejemplo de cómo hacerlo en los modelos `User` y `Task`:
```go
// internal/domain/user.go
// ...
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

// internal/domain/task.go
// ...
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
```

Es importante notar que GORM maneja automáticamente las operaciones de soft delete cuando utilizas los métodos estándar de eliminación. Por ejemplo, al llamar a `db.Delete(&user)` o `db.Delete(&task)`, GORM establecerá la marca de tiempo en el campo `DeletedAt` en lugar de eliminar físicamente el registro.

Esto significa que las consultas normales no devolverán registros marcados como eliminados. Si deseas incluir registros eliminados en tus consultas, puedes utilizar el método `Unscoped()` de GORM. Aquí tienes un ejemplo de cómo hacerlo en el repositorio:
```go
// internal/repository/user_repository.go
// ...
// GetAllIncludingDeleted obtiene todos los usuarios, incluyendo los eliminados
func (r *userRepository) GetAllIncludingDeleted() ([]domain.User, error) {
    var users []domain.User 
    err := r.db.Unscoped().Find(&users).Error
    return users, err
}
```

Este enfoque te permite aprovechar la funcionalidad de soft delete de GORM para gestionar tus datos de manera más segura y flexible.

---
## Configurar logger y connection pool en GORM
Es importante configurar adecuadamente el logger y el connection pool en GORM para optimizar el rendimiento y la depuración de tu aplicación. 

Según el modo de la aplicación (desarrollo o producción), puedes ajustar el nivel de logging para obtener más o menos detalles en los logs. Además, configurar el connection pool ayuda a gestionar las conexiones a la base de datos de manera eficiente, evitando sobrecargas y mejorando el rendimiento.

```go
// Configurar logger de GORM según el modo
var gormLogger logger.Interface
if AppConfig.GinMode == "debug" {
	gormLogger = logger.Default.LogMode(logger.Info)
} else {
	gormLogger = logger.Default.LogMode(logger.Silent)
}

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	Logger: gormLogger,
})
if err != nil {
	return fmt.Errorf("error al conectar a la base de datos: %w", err)
}
```

**Configurar connection pool**:
```go

// Configurar connection pool
sqlDB, err := db.DB()
if err != nil {
	return fmt.Errorf("error al obtener instancia SQL: %w", err)
}

sqlDB.SetMaxIdleConns(10)              // Conexiones inactivas en el pool
sqlDB.SetMaxOpenConns(100)             // Máximo de conexiones abiertas
sqlDB.SetConnMaxLifetime(time.Hour)    // Tiempo máximo de vida de una conexión
sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Tiempo máximo de inactividad de una conexión
```
Este código asegura que tu aplicación esté optimizada para diferentes entornos y maneje las conexiones a la base de datos de manera eficiente. Asegúrate de ajustar los valores del connection pool según las necesidades específicas de tu aplicación y la carga esperada.

- `SetMaxIdleConns(10)` : Establece el número máximo de conexiones inactivas que se mantendrán en el pool. Esto ayuda a reducir la latencia al reutilizar conexiones existentes.
- `SetMaxOpenConns(100)` : Define el número máximo de conexiones abiertas que se pueden tener simultáneamente. Esto previene que la aplicación abra demasiadas conexiones y sobrecargue la base de datos.
- `SetConnMaxLifetime(time.Hour)` : Especifica el tiempo máximo que una conexión puede estar abierta antes de ser cerrada y reemplazada. Esto ayuda a evitar problemas con conexiones que pueden volverse inestables con el tiempo.
- `SetConnMaxIdleTime(30 * time.Minute)` : Define el tiempo máximo que una conexión puede estar inactiva antes de ser cerrada. Esto ayuda a liberar recursos cuando las conexiones no se están utilizando.

Con estas configuraciones, tu aplicación estará mejor preparada para manejar la carga de trabajo y facilitará la depuración en caso de problemas con la base de datos.

---
## Conclusión
En esta sección, hemos explorado cómo utilizar GORM para interactuar con bases de datos relacionales en Go. Hemos cubierto desde la configuración inicial y la conexión a MySQL, hasta la definición de modelos, migraciones automáticas, relaciones entre entidades, operaciones CRUD, consultas avanzadas, y la configuración del logger y connection pool. Con estos conocimientos, estarás bien equipado para desarrollar aplicaciones robustas y eficientes utilizando GORM como tu ORM en Go. ¡Feliz codificación!
