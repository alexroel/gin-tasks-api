package domain

import "gorm.io/gorm"

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
