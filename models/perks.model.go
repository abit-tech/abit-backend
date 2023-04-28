package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Perk struct {
	ID          uuid.UUID `gorm:"type:uuid;uniqueIndex;primary_key"`
	Description string    `gorm:"type:varchar(1000); not null"`
	VideoID     uuid.UUID `gorm:"type:uuid; not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}
