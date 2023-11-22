package auth_service

import (
	models "darkness-awakens/db/models"
	"darkness-awakens/entity"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gookit/validate"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(entity.LoggingInUser, *gorm.DB)
	Register(entity.RegisteringUser, *gorm.DB)
}

func Login(ctx *gin.Context, db *gorm.DB) {
	var user entity.LoggingInUser
	ctx.BindJSON(&user)
	v := validate.Struct(user)
	if !v.Validate() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": v.Errors.One(),
		})
		return
	}

	var existingUser models.User
	result := db.First(&existingUser, "LOWER(email) = ?", strings.ToLower(user.Email))
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No account found with given email",
		})
		return
	}
	isPasswordCorrect := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))

	if isPasswordCorrect != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Incorrect Password",
		})
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(525600 * time.Minute) // 1 year
	claims["id"] = existingUser.Id
	claims["email"] = existingUser.Email
	claims["name"] = existingUser.DisplayName

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    existingUser.Id,
			"email": existingUser.Email,
			"name":  existingUser.DisplayName,
		},
	})

}

func Register(ctx *gin.Context, db *gorm.DB) {
	var user entity.RegisteringUser
	ctx.BindJSON(&user)
	v := validate.Struct(user)
	if !v.Validate() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": v.Errors.One(),
		})
		return
	}

	var existingUser models.User
	result := db.First(&existingUser, "LOWER(email) = ?", strings.ToLower(user.Email))
	if result.RowsAffected == 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "An account already exists with provided email.",
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	newUser := models.User{
		Email:       user.Email,
		Password:    string(hashedPassword),
		DisplayName: user.DisplayName,
	}
	createResult := db.Create(&newUser)
	if createResult.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(525600 * time.Minute) // 1 year
	claims["id"] = newUser.Id
	claims["email"] = newUser.Email
	claims["name"] = newUser.DisplayName

	tokenString, err := token.SignedString([]byte(string(os.Getenv("JWT_SECRET"))))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    newUser.Id,
			"email": newUser.Email,
			"name":  newUser.DisplayName,
		},
	})

}
