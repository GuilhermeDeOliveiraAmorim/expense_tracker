package util

import (
	"errors"
	"net/http"
)

type ProblemDetailsOutputDTO struct {
	ProblemDetails []ProblemDetails `json:"problem_details"`
}

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

func NewProblemDetails(t string, title string, status int, detail string, instance string) (*ProblemDetails, error) {
	pd := ProblemDetails{
		Type:     t,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}

	if err := pd.Validate(); err != nil {
		return nil, err
	}

	return &pd, nil
}

func (pd *ProblemDetails) Validate() error {
	if pd.Type == "" || len(pd.Type) > 100 {
		NewLoggerError(
			http.StatusBadRequest,
			"O tipo deve ser não vazio e ter no máximo 100 caracteres",
			"NewProblemDetails",
			"Entities",
			"Error",
		)

		return errors.New("type deve ser não vazio e ter no máximo 100 caracteres")
	}

	if pd.Title == "" || len(pd.Title) > 100 {
		NewLoggerError(
			http.StatusBadRequest,
			"O título deve ser não vazio e ter no máximo 100 caracteres",
			"NewProblemDetails",
			"Entities",
			"Error",
		)

		return errors.New("title deve ser não vazio e ter no máximo 100 caracteres")
	}

	if pd.Status < 100 || pd.Status >= 600 {
		NewLoggerError(
			http.StatusBadRequest,
			"O status deve ser um código HTTP válido",
			"NewProblemDetails",
			"Entities",
			"Error",
		)

		return errors.New("status deve ser um código HTTP válido")
	}

	if pd.Detail == "" || len(pd.Detail) > 255 {
		NewLoggerError(
			http.StatusBadRequest,
			"O detalhe deve ser não vazio e ter no máximo 255 caracteres",
			"NewProblemDetails",
			"Entities",
			"Error",
		)

		return errors.New("detail deve ser não vazio e ter no máximo 255 caracteres")
	}

	if len(pd.Instance) > 255 {
		NewLoggerError(
			http.StatusBadRequest,
			"A instância não deve ter mais do que 255 caracteres",
			"NewProblemDetails",
			"Entities",
			"Error",
		)

		return errors.New("instance não deve ter mais do que 255 caracteres")
	}

	return nil
}
