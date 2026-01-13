package domain

import "gorm.io/gorm"

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
