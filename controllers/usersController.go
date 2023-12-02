package controllers

import (
	"errors"
	"go-jwt/initializers"
	"go-jwt/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User already exists",
			})
		} else {
			// Trata outros tipos de erros
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
		}
	}

	userResponse := models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    userResponse,
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user models.User

	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid credentials",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid credentials",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte("teste123"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successfully",
	})
}
