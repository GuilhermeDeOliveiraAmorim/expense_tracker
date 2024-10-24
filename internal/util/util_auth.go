package util

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "Unauthorized",
				Title:    "Missing Authorization Header",
				Status:   http.StatusUnauthorized,
				Detail:   "Authorization header is required",
				Instance: RFC401,
			}})
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")
		if len(tokenString) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "Unauthorized",
				Title:    "Invalid Authorization Format",
				Status:   http.StatusUnauthorized,
				Detail:   "Authorization header must be in the format 'Bearer <token>'",
				Instance: RFC401,
			}})
			return
		}

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.SECRETS_LOCAL.JWT_SECRET), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "Unauthorized",
				Title:    "Invalid Token",
				Status:   http.StatusUnauthorized,
				Detail:   "Token could not be parsed or is invalid",
				Instance: RFC401,
			}})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "Unauthorized",
				Title:    "Invalid Token",
				Status:   http.StatusUnauthorized,
				Detail:   "Token is not valid",
				Instance: RFC401,
			}})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)
		c.Set("userID", userID)
		c.Next()
	}
}
