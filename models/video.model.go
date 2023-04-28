package models

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID          uuid.UUID `gorm:"type:uuid;uniqueIndex;primary_key;"`
	Name        string    `gorm:"type:varchar(100);not null"`
	CreatorID   uuid.UUID `gorm:"type:uuid;not null"`
	Description string    `gorm:"type:varchar(1000); not null"`
	ReleaseDate time.Time `gorm:"not null"`
	Status      string    `gorm:"not null"`

	TokenIcon string `gorm:"default:'default_video_icon.png';"`
	AlbumArt  string `gorm:"default:'default_video_album.png';"`

	TrailerLink string `gorm:"type:varchar(500)"`
	VideoLink   string `gorm:"type:varchar(500)"`
	HypeLink    string `gorm:"type:varchar(500)"`

	TokensOffered int32   `gorm:"type:int"`
	TokenPrice    float32 `gorm:"type:float"`
	TokensSold    int32   `gorm:"type:int"`
	RevenueShared float32 `gorm:"type:float"`

	RevenueSharingContractLink string `gorm:"type:varchar(500)"`
	OwnershipContractLink      string `gorm:"type:varchar(500)"`

	YoutubeViews    int32 `gorm:"type:int"`
	YoutubeRevenue  int32 `gorm:"type:int"`
	AudienceRevenue int32 `gorm:"type:int"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (video *Video) FractionalizeVideo() ([]Token, error) {
	if video.TokensOffered == 0 {
		return nil, errors.New("zero token count")
	}
	tokens := make([]Token, video.TokensOffered)
	for i := 0; i < int(video.TokensOffered); i++ {
		token := Token{
			Number:  i + 1,
			VideoID: video.ID.String(),
			Price:   video.TokenPrice,
		}
		tokens[i] = token
	}
	return tokens, nil
}

func (video *Video) BeforeCreate(*gorm.DB) error {
	video.ID = uuid.NewV4()
	return nil
}

type CreateVideoRequest struct {
	Name                       string   `json:"name"`
	Description                string   `json:"description"`
	ReleaseDate                string   `json:"releaseDate"`
	Perks                      []string `json:"perks"`
	TokenIcon                  string   `json:"tokenIcon"`
	AlbumArt                   string   `json:"albumArt"`
	TrailerLink                string   `json:"trailerLink"`
	VideoLink                  string   `json:"videoLink"`
	TokenPrice                 float32  `json:"tokenPrice"`
	HypeLink                   string   `json:"hypeLink"`
	TokensOffered              int32    `json:"tokensOffered"`
	RevenueShared              float32  `json:"revenueShared"`
	RevenueSharingContractLink string   `json:"revenueSharingContractLink"`
	OwnershipContractLink      string   `json:"ownershipContractLink"`
}

type CreateVideoResponse struct {
	ID            string   `json:"id"`
	CreatorID     string   `json:"creatorID"`
	Name          string   `json:"string"`
	Status        string   `json:"status"`
	Description   string   `json:"description"`
	ReleaseDate   string   `json:"releaseDate"`
	TokenIcon     string   `json:"tokenIcon"`
	Perks         []string `json:"perks"`
	AlbumArt      string   `json:"albumArt"`
	TrailerLink   string   `json:"trailerLink"`
	VideoLink     string   `json:"videoLink"`
	HypeLink      string   `json:"hypeLink"`
	TokensOffered int32    `json:"tokensOffered"`
	RevenueShared float32  `json:"revenueShared"`
}

// do not allow update of tokens offered or revenue shared

type UpdateVideoRequest struct {
	ID                         string   `json:"id"`
	Name                       string   `json:"name"`
	Description                string   `json:"description"`
	ReleaseDate                string   `json:"releaseDate"`
	TokenIcon                  string   `json:"tokenIcon"`
	AlbumArt                   string   `json:"albumArt"`
	TrailerLink                string   `json:"trailerLink"`
	Perks                      []string `json:"perks"`
	VideoLink                  string   `json:"videoLink"`
	HypeLink                   string   `json:"hypeLink"`
	RevenueSharingContractLink string   `json:"revenueSharingContractLink"`
	OwnershipContractLink      string   `json:"ownershipContractLink"`
}

type UpdateVideoResponse struct {
	Name                       string   `json:"name"`
	Description                string   `json:"description"`
	ReleaseDate                string   `json:"releaseDate"`
	Perks                      []string `json:"perks"`
	TokenIcon                  string   `json:"tokenIcon"`
	AlbumArt                   string   `json:"albumArt"`
	TrailerLink                string   `json:"trailerLink"`
	VideoLink                  string   `json:"videoLink"`
	HypeLink                   string   `json:"hypeLink"`
	RevenueSharingContractLink string   `json:"revenueSharingContractLink"`
	OwnershipContractLink      string   `json:"ownershipContractLink"`
}

type GetVideoRequest struct {
	ID string `json:"id"`
}

type GetVideoResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	CreatorID   string   `json:"creatorID"`
	Description string   `json:"description"`
	ReleaseDate string   `json:"releaseDate"`
	Perks       []string `json:"perks"`
	Status      string   `gorm:"not null"`

	TokenIcon string `json:"tokenIcon"`
	AlbumArt  string `json:"albumArt"`

	TrailerLink string `gorm:"type:varchar(500)"`
	VideoLink   string `gorm:"type:varchar(500)"`
	HypeLink    string `gorm:"type:varchar(500)"`

	TokensOffered int32   `json:"tokensOffered"`
	TokensSold    int32   `json:"tokensSold"`
	RevenueShared float32 `json:"revenueShared"`

	RevenueSharingContractLink string `json:"revenueSharingContractLink"`
	OwnershipContractLink      string `json:"ownershipContractLink"`

	YoutubeViews    int32 `json:"youtubeViews"`
	YoutubeRevenue  int32 `json:"youtubeRevenue"`
	AudienceRevenue int32 `json:"audienceRevenue"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteVideoRequest struct {
	ID string `json:"id"`
}

type DeleteVideoResponse struct {
	Success bool `json:"success"`
}

type FetchVideosForCreatorDashRequest struct {
	CreatorID string `json:"creatorID"`
}

type FetchVideosForCreatorDashResponse struct {
	CreatorID string         `json:"creatorID"`
	Videos    []VideoForDash `json:"videos"`
}

type VideoForDash struct {
	Name          string  `json:"string"`
	Status        string  `json:"status"`
	ReleaseDate   string  `json:"releaseDate"`
	TokenIcon     string  `json:"tokenIcon"`
	AlbumArt      string  `json:"albumArt"`
	TokensOffered int32   `json:"tokensOffered"`
	RevenueShared float32 `json:"revenueShared"`
}
