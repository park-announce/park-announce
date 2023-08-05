package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/entity"
)

func UseUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if isAuthorizationRequired(c.Request.URL.Path) && isAuthorizationEnabled() {
			authorization := c.Request.Header.Get("Authorization")
			if authorization == "" {
				authorization = c.Request.URL.Query().Get("Authorization")
			}

			// get user-id from client in http request header with key "Authorization" in JWT format
			if authorization == "" {
				c.AbortWithError(http.StatusUnauthorized, errors.New("no authorization header found"))
			}

			user := entity.User{}
			token, err := jwt.ParseWithClaims(authorization, &user, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("PA_API_JWT_KEY")), nil
			})

			if err != nil {
				log.Println("error :", err)
				c.AbortWithError(http.StatusUnauthorized, errors.New("invalid authorization"))
			}

			if token.Valid {
				c.Set("User", user)
			} else {
				c.AbortWithError(http.StatusUnauthorized, errors.New("invalid authorization"))
			}
		}

		c.Next()
	}
}

func isAuthorizationRequired(path string) bool {
	unAuthoriziedPaths := []string{"/preregisterations", "/registerationverifications", "/registerations", "/preregisterations/google", "/google/oauth2/token", "/corporation/token"}
	var isAuthorizationRequired bool = true
	if len(unAuthoriziedPaths) > 0 && path != "" {
		for _, unAuthoriziedPath := range unAuthoriziedPaths {
			if unAuthoriziedPath == path {
				isAuthorizationRequired = false
				break
			}
		}
	}

	return isAuthorizationRequired
}

func isAuthorizationEnabled() bool {
	return true
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
