package util

import "regexp"

func IsValidHexColor(hexColor string) bool {
	regex := regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
	return regex.MatchString(hexColor)
}
