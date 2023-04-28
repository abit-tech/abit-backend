package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"www.github.com/abit-tech/abit-backend/common"
	"www.github.com/abit-tech/abit-backend/htmltopdf"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
)

func CreateVideo(ctx *gin.Context) {
	var payload *models.CreateVideoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.Role == common.RoleUser {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "only creators can create videos",
		})
		return
	}

	now := time.Now()
	id := uuid.NewV4()

	// todo decide on the format of date that will be passed by rahul
	relDate, err := time.Parse("Jan 02, 2006", payload.ReleaseDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "time format should be YYYY-MM-DD",
		})
		return
	}

	newVideo := models.Video{
		ID:          id,
		Name:        payload.Name,
		CreatorID:   currentUser.ID,
		Description: payload.Description,
		ReleaseDate: relDate,
		Status:      common.VideoStatusPending,

		TokenIcon: payload.TokenIcon,
		AlbumArt:  payload.AlbumArt,

		TrailerLink: payload.TrailerLink,
		VideoLink:   payload.VideoLink,
		HypeLink:    payload.HypeLink,

		TokensOffered: payload.TokensOffered,
		TokenPrice:    payload.TokenPrice,
		TokensSold:    0,
		RevenueShared: payload.RevenueShared,

		RevenueSharingContractLink: payload.RevenueSharingContractLink,
		OwnershipContractLink:      payload.OwnershipContractLink,

		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newVideo)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "something went wrong",
		})
		return
	}

	// create tokens for this video
	go FractionalizeVideo(newVideo)

	// create revenue sharing contract
	go htmltopdf.GenerateRevenueSharingContract(newVideo, currentUser)

	// create perks for this video
	for _, data := range payload.Perks {
		perkID := uuid.NewV4()
		newPerk := models.Perk{
			ID:          perkID,
			Description: data,
			VideoID:     newVideo.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		tx := initializers.DB.Create(&newPerk)
		if tx.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "error",
				"message": "something went wrong while creating perks",
			})
			return
		}
	}

	// add duplicate check for video URL, cannot be same

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"video": newVideo,
		},
	})

}

func FractionalizeVideo(video models.Video) {
	tokens, err := video.FractionalizeVideo()
	if err != nil {
		fmt.Printf("error in creating tokens: %v\n", err.Error())
		return
	}

	result := initializers.DB.Create(&tokens)
	if result.Error != nil {
		fmt.Printf("error in inserting tokens: %v\n", result.Error.Error())
		return
	}

	fmt.Print("successfully created tokens")
}

func UpdateVideo(ctx *gin.Context) {
	var payload *models.UpdateVideoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	video := models.Video{}
	result := initializers.DB.First(&video, "id = ?", strings.ToLower(payload.ID))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid video id",
		})
		return
	}

	if video.CreatorID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "not authorized",
		})
		return
	}

	now := time.Now()

	// change this approach so that newVideo is initially same as old video and then
	// update the new fields from request
	newVideo := models.Video{
		ID:              video.ID,
		CreatorID:       video.CreatorID,
		UpdatedAt:       now,
		Status:          video.Status,
		TokensOffered:   video.TokensOffered,
		TokensSold:      video.TokensSold,
		RevenueShared:   video.RevenueShared,
		YoutubeViews:    video.YoutubeViews,
		YoutubeRevenue:  video.YoutubeRevenue,
		AudienceRevenue: video.AudienceRevenue,
		CreatedAt:       video.CreatedAt,
	}

	if payload.Name != "" {
		newVideo.Name = payload.Name
	} else {
		newVideo.Name = video.Name
	}

	if payload.Description != "" {
		newVideo.Description = payload.Description
	} else {
		newVideo.Description = video.Description
	}

	if payload.ReleaseDate != "" {
		relDate, err := time.Parse("Jan 02, 2006", payload.ReleaseDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "time format should be YYYY-MM-DD",
			})
			return
		}
		newVideo.ReleaseDate = relDate
	} else {
		newVideo.ReleaseDate = video.ReleaseDate
	}

	if payload.TokenIcon != "" {
		newVideo.TokenIcon = payload.TokenIcon
	} else {
		newVideo.TokenIcon = video.TokenIcon
	}

	if payload.AlbumArt != "" {
		newVideo.AlbumArt = payload.AlbumArt
	} else {
		newVideo.AlbumArt = video.AlbumArt
	}

	if payload.TrailerLink != "" {
		newVideo.TrailerLink = payload.TrailerLink
	} else {
		newVideo.TrailerLink = video.TrailerLink
	}

	if payload.VideoLink != "" {
		newVideo.VideoLink = payload.VideoLink
	} else {
		newVideo.VideoLink = video.VideoLink
	}

	if payload.HypeLink != "" {
		newVideo.HypeLink = payload.HypeLink
	} else {
		newVideo.HypeLink = video.HypeLink
	}

	if payload.Description != "" {
		newVideo.Description = payload.Description
	} else {
		newVideo.Description = video.Description
	}

	if payload.RevenueSharingContractLink != "" {
		newVideo.RevenueSharingContractLink = payload.RevenueSharingContractLink
	} else {
		newVideo.RevenueSharingContractLink = video.RevenueSharingContractLink
	}

	if payload.OwnershipContractLink != "" {
		newVideo.OwnershipContractLink = payload.OwnershipContractLink
	} else {
		newVideo.OwnershipContractLink = video.OwnershipContractLink
	}

	tx := initializers.DB.Model(&newVideo).Where("id = ?", video.ID).UpdateColumns(&newVideo)
	if tx.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"video": newVideo,
		},
	})
}

func GetVideo(ctx *gin.Context) {
	var payload *models.GetVideoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	video := models.Video{}
	result := initializers.DB.First(&video, "id = ?", strings.ToLower(payload.ID))

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "video not found",
		})
		return
	}

	// fetch all the perks for a video
	perks := []models.Perk{}
	tx := initializers.DB.Find(&perks).Where("video_id = ?", strings.ToLower(video.ID.String()))
	if tx.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "error in fetching perks",
		})
		return
	}
	perkString := make([]string, 0)
	for _, perk := range perks {
		perkString = append(perkString, perk.Description)
	}

	resp := models.GetVideoResponse{
		ID:          video.ID.String(),
		Name:        video.Name,
		CreatorID:   video.CreatorID.String(),
		Description: video.Description,
		ReleaseDate: video.ReleaseDate.String(),
		Status:      video.Status,
		TokenIcon:   video.TokenIcon,
		AlbumArt:    video.AlbumArt,
		TrailerLink: video.TrailerLink,
		VideoLink:   video.VideoLink,
		HypeLink:    video.HypeLink,

		TokensOffered: video.TokensOffered,
		TokensSold:    video.TokensSold,
		RevenueShared: video.RevenueShared,

		RevenueSharingContractLink: video.RevenueSharingContractLink,
		OwnershipContractLink:      video.OwnershipContractLink,

		CreatedAt: video.CreatedAt,
		UpdatedAt: video.UpdatedAt,

		Perks: perkString,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"video": resp,
		},
	})
}

func DeleteVideo(ctx *gin.Context) {
	var payload *models.DeleteVideoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	video := models.Video{}
	result := initializers.DB.First(&video, "id = ?", strings.ToLower(payload.ID))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid video id",
		})
		return
	}

	if video.CreatorID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "not authorized",
		})
		return
	}

	tx := initializers.DB.Delete(&video, "id = ?", strings.ToLower(payload.ID))
	if tx.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	// delete the perks
	perk := models.Perk{}
	tx = initializers.DB.Where("video_id = ?", strings.ToLower(video.ID.String())).Delete(&perk)
	if tx.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"success": true,
		},
	})

}

func PurchaseToken(ctx *gin.Context) {
	var payload *models.PurchaseTokenRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)

	video := models.Video{}
	result := initializers.DB.First(&video, "id = ?", strings.ToLower(payload.VideoID))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid video id",
		})
		return
	}

	if currentUser.ID.String() == video.CreatorID.String() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "cannot purchase token of your own video",
		})
		return
	}

	token := models.Token{}
	// check if this user already has a token of this video
	tx := initializers.DB.Limit(1).Order("number").Find(&token, "video_id = ? AND owner_id = ?",
		strings.ToLower(payload.VideoID),
		strings.ToLower(currentUser.ID.String()),
	)
	// if you find something, throw error
	if tx.Error != nil {
		fmt.Printf("something went wrong: %v\n", tx.Error.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	if token.ID != uuid.Nil {
		fmt.Printf("already owned token found \n")
		fmt.Printf("token found: %v\n", token)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "cannot purchase multiple tokens of the same video",
		})
		return
	}

	nilUUID := uuid.Nil
	tx = initializers.DB.Limit(1).Order("number").Find(&token, "video_id = ? AND owner_id = ?",
		strings.ToLower(payload.VideoID),
		strings.ToLower(string(nilUUID.String())),
	)

	if tx.Error != nil {
		fmt.Printf("error in fetching token: %v\n", tx.Error.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "could not find token",
		})
		return
	}

	creator := models.User{}
	tx = initializers.DB.First(&creator, "id = ?", strings.ToLower(video.CreatorID.String()))
	if tx.Error != nil {
		fmt.Printf("error in fetching creator: %v\n", tx.Error.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "could not purchase token",
		})
		return
	}
	go htmltopdf.GenerateTokenOwnershipContract(token, video, creator, currentUser)

	fmt.Printf("token found: %v\n", token)
	token.OwnerID = currentUser.ID
	tx = initializers.DB.Model(&token).Where("id = ?", token.ID).UpdateColumns(&token)
	if tx.Error != nil {
		fmt.Printf("error in updating token: %v\n", tx.Error.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "could not purchase token",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
	})
}
