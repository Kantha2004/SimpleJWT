package models

type User struct {
	Username     string `json:"username" gorm:"unique;not null"`
	Email        string `json:"email" gorm:"unique;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	TableModel
}

type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
