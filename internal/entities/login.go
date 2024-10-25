package entities

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/config"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLogin(email, password string) (*Login, []util.ProblemDetails) {
	validationErrors := ValidateLogin(email, password)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Login{
		Email:    email,
		Password: password,
	}, nil
}

func ValidateLogin(email, password string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if !isValidEmail(email) {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Invalid email",
			Status:   400,
			Detail:   "Email is invalid",
			Instance: util.RFC400,
		})
	}

	if !isValidPassword(password) {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Invalid password",
			Status:   400,
			Detail:   "Password must be at least 6 characters long, contain at least one uppercase letter, one lowercase letter, one digit, and one special character",
			Instance: util.RFC400,
		})
	}

	return validationErrors
}

func isValidEmail(email string) bool {
	emailPattern := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	match, _ := regexp.MatchString(emailPattern, email)
	return match
}

func isValidPassword(password string) bool {
	return hasMinimumLength(password, 6) &&
		hasUpperCaseLetter(password) &&
		hasLowerCaseLetter(password) &&
		hasDigit(password) &&
		hasSpecialCharacter(password)
}

func hasMinimumLength(password string, length int) bool {
	return len(password) >= length
}

func hasUpperCaseLetter(password string) bool {
	return strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func hasLowerCaseLetter(password string) bool {
	return strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
}

func hasDigit(password string) bool {
	return strings.ContainsAny(password, "0123456789")
}

func hasSpecialCharacter(password string) bool {
	specialCharacters := "@#$%&*"
	return strings.ContainsAny(password, specialCharacters)
}

func hashString(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (lo *Login) CompareAndDecrypt(hashedData string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(data))
	return err == nil
}

func (lo *Login) EncryptEmail() error {
	key := []byte(config.SECRETS_VAR.JWT_SECRET)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(lo.Email))
	lo.Email = hex.EncodeToString(h.Sum(nil))

	return nil
}

func (lo *Login) EncryptPassword() error {
	hashedPassword, err := hashString(lo.Password)
	if err != nil {
		return err
	}

	lo.Password = string(hashedPassword)

	return nil
}

func (lo *Login) DecryptEmail(email string) bool {
	if lo.CompareAndDecrypt(lo.Email, email) {
		return true
	} else {
		return false
	}
}

func (lo *Login) DecryptPassword(password string) bool {
	if lo.CompareAndDecrypt(lo.Password, password) {
		return true
	} else {
		return false
	}
}

func (lo *Login) ChangeEmail(newEmail string) {
	lo.Email = newEmail
}

func (lo *Login) ChangePassword(newPassword string) {
	lo.Password = newPassword
}

func (lo *Login) Equals(other *Login) bool {
	return lo.Email == other.Email && lo.Password == other.Password
}
