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
			c.JSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7235#section-3.1",
				Title:    "Unauthorized",
				Status:   http.StatusUnauthorized,
				Detail:   "Missing or invalid Authorization header",
				Instance: c.Request.URL.String(),
			}})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")
		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7235#section-3.1",
				Title:    "Unauthorized",
				Status:   http.StatusUnauthorized,
				Detail:   "Invalid Authorization header format",
				Instance: c.Request.URL.String(),
			}})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.SECRETS_LOCAL.JWT_SECRET), nil
		})

		if err != nil {
			fmt.Println("JWT Parse Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7235#section-3.1",
				Title:    "Unauthorized",
				Status:   http.StatusUnauthorized,
				Detail:   "Invalid token",
				Instance: c.Request.URL.String(),
			}})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Println("Invalid Token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7235#section-3.1",
				Title:    "Unauthorized",
				Status:   http.StatusUnauthorized,
				Detail:   "Token is not valid",
				Instance: c.Request.URL.String(),
			}})
			c.Abort()
			return
		}

		c.Set("user", token.Claims.(jwt.MapClaims))
		c.Next()
	}
}
