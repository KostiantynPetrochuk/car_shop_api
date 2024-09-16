package routes

import (
	"fmt"
	"net/http"
	"time"

	"example.com/models"
	"example.com/utils"
	"github.com/gin-gonic/gin"
)

const ACCESS_TOKEN_EXPIRE_TIME = time.Hour * 1

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.Save()
	if err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func signin(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticated user."})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.Login, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticated user."})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(user.Login, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticated user."})
		return
	}

	now := time.Now()
	expiresAt := now.Add(ACCESS_TOKEN_EXPIRE_TIME)
	expiresIn := expiresAt.UnixNano() / int64(time.Millisecond)
	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful.",
		"user":    user,
		"tokens": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"expiresIn":    expiresIn,
		}})
}

func refresh(context *gin.Context) {
	userId, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context."})
		return
	}
	userIdInt, ok := userId.(int64)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "User ID is not of type int64."})
		return
	}
	login, exists := context.Get("login")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context."})
		return
	}
	loginStr, ok := login.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Login is not of type string."})
		return
	}
	accessToken, err := utils.GenerateAccessToken(loginStr, userIdInt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticated user."})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(loginStr, userIdInt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticated user."})
		return
	}
	now := time.Now()
	expiresAt := now.Add(ACCESS_TOKEN_EXPIRE_TIME)
	expiresIn := expiresAt.UnixNano() / int64(time.Millisecond)

	context.JSON(http.StatusOK, gin.H{"message": "Login successful.", "tokens": gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expiresIn":    expiresIn,
	}})
}
