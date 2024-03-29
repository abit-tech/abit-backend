package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Creator struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name     string    `gorm:"type:varchar(100);not null"`
	Email    string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string    `gorm:"not null"`

	ChannelLink string `gorm:"type:varchar(200);not null"`
	Photo       string `gorm:"default:'default_creator_pic.png';"`

	Verified  bool      `gorm:"default:false;"`
	Provider  string    `gorm:"default:'local';"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (creator *Creator) BeforeCreate(*gorm.DB) error {
	creator.ID = uuid.NewV4()
	return nil
}

type RegisterCreatorInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginCreatorInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreatorResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"Name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
	Provider  string `json:"provider,omitempty"`
	Photo     string `json:"photo,omitempty"`
	Verified  bool   `json:"verified,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FilteredResponse is used to omit sensitive fields from GORM response
// func FilteredResponse(user *User) UserResponse {
// 	return UserResponse{
// 		ID:        user.ID.String(),
// 		Name:      user.Name,
// 		Email:     user.Email,
// 		Role:      user.Role,
// 		Verified:  user.Verified,
// 		Photo:     user.Photo,
// 		Provider:  user.Provider,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	}
// }
