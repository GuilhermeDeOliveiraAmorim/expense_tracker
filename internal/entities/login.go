package entities

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLogin(email, password string) (*Login, []error) {
	validationErrors := ValidateLogin(email, password)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Login{
		Email:    email,
		Password: password,
	}, nil
}

func ValidateLogin(email, password string) []error {
	var validationErrors []error

	if !isValidEmail(email) {
		validationErrors = append(validationErrors, fmt.Errorf("invalid email"))
	}

	if !isValidPassword(password) {
		validationErrors = append(validationErrors, fmt.Errorf("invalid password"))
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

func compareAndDecrypt(hashedData string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(data))
	return err == nil
}

func (lo *Login) EncryptEmail() (string, error) {
	hashedEmail, err := hashString(lo.Email)
	if err != nil {
		return "", err
	}
	return hashedEmail, nil
}

func (lo *Login) EncryptPassword() (string, error) {
	hashedPassword, err := hashString(lo.Password)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func (lo *Login) DecryptEmail(email string) bool {
	if compareAndDecrypt(lo.Email, email) {
		return true
	} else {
		return false
	}
}

func (lo *Login) DecryptPassword(password string) bool {
	if compareAndDecrypt(lo.Password, password) {
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
