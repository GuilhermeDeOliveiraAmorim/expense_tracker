package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (string, []util.ProblemDetails) {
	var problemsDetails []util.ProblemDetails

	userID, exists := c.Get("userID")
	if !exists {
		return "", append(problemsDetails, util.ProblemDetails{
			Type:     "Unauthorized",
			Title:    "Missing User ID",
			Status:   http.StatusUnauthorized,
			Detail:   "User id is required",
			Instance: util.RFC401,
		})
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		return "", append(problemsDetails, util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Invalid User ID",
			Status:   http.StatusBadRequest,
			Detail:   "A valid user id is required",
			Instance: util.RFC400,
		})
	}

	return userIDStr, problemsDetails
}
