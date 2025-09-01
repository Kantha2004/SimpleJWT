package handlers

import (
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/gin-gonic/gin"
)

func (d *Dependencies) CreateUser(c *gin.Context) {
	var req models.CreateUser

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userRepo := db.NewUserRepository(d.DB)

	userNameExists, err := userRepo.UserNameExists(req.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if userNameExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	emailExists, err := userRepo.EmailExists(req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if emailExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating the user"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	userId, err := userRepo.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating the user"})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created succefully!!!",
		"userId":  userId,
	})
}
