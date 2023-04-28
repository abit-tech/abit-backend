package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Email       string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password    string    `gorm:"not null"`
	Role        string    `gorm:"type:varchar(20);default:'user';"` // can be creator
	ChannelLink string    `gorm:"type:varchar(200);"`
	Photo       string    `gorm:"default:'default.png';"`
	Verified    bool      `gorm:"default:false;"`
	Provider    string    `gorm:"default:'manual';"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (user *User) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.NewV4()
	return nil
}

type RegisterUserInput struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Role        string `json:"role" binding:"required"`
	ChannelLink string `json:"channelLink"`
	Password    string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"Name,omitempty"`
	Email       string    `json:"email,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	Photo       string    `json:"photo,omitempty"`
	ChannelLink string    `json:"channelLink"`
	Verified    bool      `json:"verified,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type GetPurchasedTokensForUserRequest struct {
	UserID string `json:"userID"`
}

type GetPurchasedTokensForUserResponse struct {
	UserID string         `json:"userID"`
	Tokens []TokenForDash `json:"tokens"`
}

// FilteredResponse is used to omit sensitive fields from GORM response
func FilteredResponse(user *User) UserResponse {
	return UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		Verified:    user.Verified,
		ChannelLink: user.ChannelLink,
		Photo:       user.Photo,
		Provider:    user.Provider,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
