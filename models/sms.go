package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *VerificationSMS) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type (
	VerificationSMS struct {
		ID       uuid.UUID
		Pin      int
		Expire   int64
		Confirm  bool
		Attempts int
		Number   string
	}

	Pin struct {
		ID  uuid.UUID `json:"validation_id"`
		Pin int       `json:"pin"`
	}
)
