package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)


func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	//base.CreatedAt = time.Now()
	//base.UpdatedAt = time.Now()
	//base.DeletedAt = nil
	return
}


type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}