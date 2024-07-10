package controllers

import (
	"DellTradingApi/dtos"
	"DellTradingApi/models"
	"DellTradingApi/services"

	"github.com/gin-gonic/gin"

	"net/http"
)

func InitUserRoutes(router *gin.RouterGroup) {
	router.POST("/register", RegisterUser)
	router.POST("/login", LoginUser)
	router.POST("/authorize", AuthorizeUser)
	router.POST("/info", RetrieveUserInfo)
}

func RegisterUser(c *gin.Context) {
	//if request can be parsed correctly and matches all fields needed to register a user
	var json dtos.RegisterRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		//todo: more verbose json errors
		c.JSON(http.StatusBadRequest, err)
		return
	}

	//make user
	var newUser *models.UserEntity
	var err error
	if newUser, err = services.CreateUser(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, sucessfulLoginResponse(newUser))
}

// todo: should I prevent a user from loggin in twice or should I leave that to the front-end
func LoginUser(c *gin.Context) {
	//ensures that login request contains proper fields
	var json dtos.LoginRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//authenticates
	user := &models.UserEntity{
		Email: json.Email}
	user, err := services.AuthenticateWithPassword(user, json.Password)

	//handles errors
	//todo: maybe have more specific errors but maybe not because we do not want user to know why login failed
	if err != nil {
		c.JSON(http.StatusForbidden, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, sucessfulLoginResponse(user))
}

func AuthorizeUser(c *gin.Context) {
	if _, exists := c.Get("user_id"); exists {
		c.JSON(http.StatusOK, gin.H{"message": "You are authorized"})
	}
}

func RetrieveUserInfo(c *gin.Context) {
	id, _ := c.Get("user_id")
	userId, validInt := id.(uint)
	if validInt {
		user, err := services.GetUserById(userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid user"})
			return
		}

		c.JSON(http.StatusOK, userInfoResponse(user))
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid user id format"})
		return
	}
}

func sucessfulLoginResponse(user *models.UserEntity) gin.H {
	return gin.H{
		"Authorization": services.GenerateJwtToken(user),
		"user":          userInfoResponse(user),
	}
}

func userInfoResponse(user *models.UserEntity) gin.H {
	return gin.H{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"cash":       user.Cash,
	}
}
