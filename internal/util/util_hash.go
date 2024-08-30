package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/configs"
)

func HashEmailWithHMAC(email string) (string, []ProblemDetails) {
	configs, LoadConfigErr := configs.LoadConfig(".")
	if LoadConfigErr != nil {
		return "", []ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error loading configuration",
				Status:   500,
				Detail:   LoadConfigErr.Error(),
				Instance: RFC500,
			},
		}
	}

	key := []byte(configs.JwtSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(email))

	return hex.EncodeToString(h.Sum(nil)), nil
}
