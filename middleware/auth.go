package middleware

import (
	"DellTradingApi/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func EnsureAuthenticated(c *gin.Context) {
	requestJwt := c.Request.Header.Get("Authorization")

	claims := &services.DellTradingJWTClaims{}
	tkn, err := jwt.ParseWithClaims(requestJwt, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//if I get a different sigining method than what I hardcoded
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		//only abort if not going to paths that don't require authentication
		var path string = c.Request.URL.Path
		if path != "/api/v0/users/register" && path != "/api/v0/users/login" {
			//stop executing api request if user is not valid
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		return
	}

	var jwtClaims *services.DellTradingJWTClaims
	var ok bool
	if jwtClaims, ok = interface{}(tkn.Claims).(*services.DellTradingJWTClaims); !ok {
		//I dont think this error will ever happen because I hardcoded the claims to be of type DellTradingJWTClaims
		log.Fatal("unknown claims type, cannot proceed")
	}

	//set user id in context so that i have acess to it during other requests
	c.Set("user_id", jwtClaims.UserID)
}

// Parse takes the token string and a function for looking up the key. The latter is especially
// // useful if you use multiple keys for your application.  The standard is to use 'kid' in the
// // head of the token to identify which key to use, but the parsed token (head and claims) is provided
// // to the callback, providing flexibility.
// token, err := jwt.Parse(requestJwt, func(token *jwt.Token) (interface{}, error) {
// 	// Don't forget to validate the alg is what you expect:
// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 	}

// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 	return []byte(os.Getenv("JWT_SECRET")), nil
// })
// if err != nil {
// 	fmt.Println("PRINTING ERROR")
// 	log.Fatal(err)
// }
