package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/entity"
)

func UseUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.Request.Header.Get("Authorization")

		if isAuthorizationRequired(c.Request.URL.Path) && isAuthorizationEnabled() {

			if authorization == "" {
				authorization = c.Request.URL.Query().Get("Authorization")
			}

			// get user-id from client in http request header with key "Authorization" in JWT format
			if authorization == "" {
				c.AbortWithError(http.StatusUnauthorized, errors.New("No Authorization header found"))
			}

			user := entity.User{}
			token, err := jwt.ParseWithClaims(authorization, &user, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("PA_API_JWT_KEY")), nil
			})

			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid Authorization"))
			}

			if token.Valid {
				userJson, err := json.Marshal(user)
				if err != nil {
					c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid Authorization"))
				}

				fmt.Println("user ->", string(userJson))
				c.Set("User", user)
			}
		}

		c.Next()
	}
}

func isAuthorizationRequired(path string) bool {
	unAuthoriziedPaths := []string{"/signin", "/signup", "/google/oauth2/token", "/google/oauth2", "/favicon.ico", "/", "/play", "/multiplay", "/stage", "/css/", "/css/play.css", "/js", "/js/play.js", "/js/multiplay.js", "/images", "/images/happy.png", "/images/angry.png", "/images/angry2.png", "/images/history.png", "/images/star.png", "/images/timer.png", "/images/mark.png"}
	var isAuthorizationRequired bool = true
	if unAuthoriziedPaths != nil && len(unAuthoriziedPaths) > 0 && path != "" {
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
