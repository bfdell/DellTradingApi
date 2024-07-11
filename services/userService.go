package services

import (
	"DellTradingApi/dtos"
	"DellTradingApi/infra"
	"DellTradingApi/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const JWT_EXPIRATION_DAYS uint32 = 2

func CreateUser(json *dtos.RegisterRequestDto) (*models.UserEntity, error) {
	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	var newUser *models.UserEntity = &models.UserEntity{
		Password:  string(password),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Email:     json.Email,
	}
	database := infra.GetDB()
	err := database.Create(newUser).Error

	return newUser, err
}

func GetUserFromContext(c *gin.Context) (*models.UserEntity, error) {
	id, _ := c.Get("user_id")
	userId, validInt := id.(uint)
	if validInt {
		user, err := GetUserById(userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid user"})
			return nil, err
		}

		return user, nil
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid user id format"})
		return nil, fmt.Errorf("invalid user ID")
	}
}

func GetUserById(userID uint) (*models.UserEntity, error) {
	return ReadUser(&models.UserEntity{Model: gorm.Model{ID: userID}})
}

func ReadUser(data *models.UserEntity) (*models.UserEntity, error) {
	database := infra.GetDB()
	err := database.Take(data).Error
	return data, err
}

// func UpdateUser(data models.UserEntity) error {
// 	database := infra.GetDB()
// 	err := database.Save(data).Error
// 	return err
// }

// func DeleteUser(data models.UserEntity) error {
// 	database := infra.GetDB()
// 	err := database.Create(data).Error
// 	return err
// }

func AuthenticateWithPassword(user *models.UserEntity, password string) (*models.UserEntity, error) {
	//will either return the error code from reading the user from the email,
	//or the error from authenticating the password
	user, err := ReadUser(user)
	if err == nil {
		return user, IsCorrectPassword(user, password)
	}

	return user, err
}

func IsCorrectPassword(u *models.UserEntity, password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

type DellTradingJWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(user *models.UserEntity) string {
	claims := &DellTradingJWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(JWT_EXPIRATION_DAYS))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
}

// func (u *UserEntity) SetPassword(password string) error {
// 	if len(password) == 0 {
// 		return errors.New("password should not be empty")
// 	}
// 	bytePassword := []byte(password)
// 	// Make sure the second param `bcrypt generator cost` between [4, 32)
// 	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
// 	u.Password = string(passwordHash)
// 	return nil
// }

// func (user *User) BeforeSave(db *gorm.DB) (err error) {
// 	if len(user.Roles) == 0 {
// 		// role := Role{}
// 		userRole := Role{}
// 		// db.Model(&role).Where("name = ?", "ROLE_USER").First(&userRole)
// 		db.Model(&Role{}).Where("name = ?", "ROLE_USER").First(&userRole)
// 		//db.Where(&models.Role{Name: "ROLE_USER"}).Attrs(models.Role{Description: "For standard Users"}).FirstOrCreate(&userRole)
// 		user.Roles = append(user.Roles, userRole)
// 	}
// 	return
// }
