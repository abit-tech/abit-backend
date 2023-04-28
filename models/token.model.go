package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key;"`
	Number                int       `gorm:"type:int; not null;"`
	VideoID               string    `gorm:"not null"`
	OwnerID               uuid.UUID `gorm:"type:uuid"`
	Price                 float32   `gorm:"type:float"`
	Icon                  string    `gorm:"default:'default_token_icon.png';"`
	TransactionID         string    `gorm:"type:varchar(150)"`
	OwnershipContractLink string    `gorm:"varchar(150)"`
	RevenueTillMonth      float32   `gorm:"type:float"`
	CreatedAt             time.Time `gorm:"not null"`
	UpdatedAt             time.Time `gorm:"not null"`
}

// tokenID isSold SoldTo VideoID
// for all tokens with videoID = currentVideo.ID, if soldTo == userID, then return error

func (token *Token) BeforeCreate(*gorm.DB) error {
	token.ID = uuid.NewV4()
	return nil
}

type PurchaseTokenRequest struct {
	VideoID string `json:"videoID"`
}

type PurchaseTokenResponse struct {
	VideoID string `json:"videoID"`
	TokenID string `json:"tokenID"`
	Number  int    `json:"number"`
}

type TokenForDash struct {
	Number           int     `json:"number"`
	VideoName        string  `json:"videoName"`
	Icon             string  `json:"icon"`
	RevenueTillMonth float32 `json:"revenueTillMonth"`
}
