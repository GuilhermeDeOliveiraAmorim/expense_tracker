package util

import (
	"time"
)

const (
	TIMEZONE   = "America/Sao_Paulo"
	DATEFORMAT = "02012006"
)

func ParseDate(date string) (time.Time, error) {
	location, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		return time.Time{}, err
	}

	newDate, err := time.ParseInLocation(DATEFORMAT, date, location)
	if err != nil {
		return time.Time{}, err
	}

	return newDate, nil
}
