package entities

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type SharedEntity struct {
	ID            string    `json:"id"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}

func NewSharedEntity() *SharedEntity {
	timeNow := time.Now()

	return &SharedEntity{
		ID:            ulid.Make().String(),
		Active:        true,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		DeactivatedAt: timeNow,
	}
}

func (se *SharedEntity) Activate() {
	timeNow := time.Now()
	se.UpdatedAt = timeNow
	se.Active = true
}

func (se *SharedEntity) Deactivate() {
	timeNow := time.Now()
	se.DeactivatedAt = timeNow
	se.UpdatedAt = timeNow
	se.Active = false
}
