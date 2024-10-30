package util

import (
	"strconv"
	"strings"
)

func FloatToBRL(value float64) string {
	valueToString := strconv.FormatFloat(value, 'f', -1, 64)
	valueToBRL := "R$ " + strings.Replace(valueToString, ".", ",", 1)
	return valueToBRL
}
