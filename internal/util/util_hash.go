package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/config"
)

func HashEmailWithHMAC(email string) (string, []ProblemDetails) {
	key := []byte(config.SECRETS_LOCAL.JWT_SECRET)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(email))

	return hex.EncodeToString(h.Sum(nil)), nil
}
