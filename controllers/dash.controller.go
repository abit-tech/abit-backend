package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
)

func CreatorDash(ctx *gin.Context) {
	var payload *models.FetchVideosForCreatorDashRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	creatorID := payload.CreatorID
	videos := []models.Video{}
	tx := initializers.DB.Where("creator_id = ?", strings.ToLower(creatorID)).Find(&videos)

	if tx.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "error in fetching videos",
		})
		return
	}

	processedVideos, err := processRawVideosToDashVideos(videos)
	if err != nil {
		// log error and return
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "error in parsing videos",
		})
		return

	}

	resp := &models.FetchVideosForCreatorDashResponse{
		CreatorID: creatorID,
		Videos:    processedVideos,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"response": resp,
		},
	})
}

func UserDash(ctx *gin.Context) {
	var payload *models.GetPurchasedTokensForUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	userID := payload.UserID
	tokens := []models.Token{}
	fmt.Printf("fetching tokens for user: %v\n", userID)
	tx := initializers.DB.Where("owner_id = ?", strings.ToLower(userID)).Find(&tokens)
	if tx.Error != nil {
		fmt.Println("error in fetching tokens")
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "error in fetching tokens",
		})
		return
	}

	fmt.Printf("fetched %d tokens\n", len(tokens))

	processedTokens, err := processRawTokenToDashToken(tokens)
	if err != nil {
		// log error and return
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "error in parsing tokens",
		})
		return

	}

	resp := &models.GetPurchasedTokensForUserResponse{
		UserID: userID,
		Tokens: processedTokens,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"response": resp,
		},
	})
}

func processRawTokenToDashToken(tokens []models.Token) ([]models.TokenForDash, error) {
	resp := make([]models.TokenForDash, 0)
	var wg sync.WaitGroup
	videoDataMap := sync.Map{}

	for _, token := range tokens {
		wg.Add(1)
		fmt.Printf("processing token: %s", token.ID.String())
		go fetchVideoNameForToken(token, &videoDataMap, &wg)
	}

	wg.Wait()
	fmt.Printf("processing finished")

	for _, token := range tokens {
		tokenVideo, ok := videoDataMap.Load(token.VideoID)
		if !ok {
			// log error
			fmt.Printf("could not find video name for tokenID: %s", token.ID.String())
			return nil, fmt.Errorf("could not find video name for tokenID: %s", token.ID.String())
		}

		videoName := tokenVideo.(models.Video).Name
		dashToken := models.TokenForDash{
			Number:           token.Number,
			VideoName:        videoName,
			Icon:             token.Icon,
			RevenueTillMonth: token.RevenueTillMonth,
		}

		resp = append(resp, dashToken)
	}

	return resp, nil
}

func fetchVideoNameForToken(token models.Token, videoDataMap *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	_, found := videoDataMap.Load(token.VideoID)
	if found {
		fmt.Printf("already found video : %s\n", token.VideoID)
		return
	}
	video := models.Video{}
	result := initializers.DB.First(&video, "id = ?", strings.ToLower(token.VideoID))
	if result.Error != nil {
		fmt.Printf("could not find video: %s\n", token.VideoID)
		return
	}
	videoDataMap.Store(token.VideoID, video)
}

func processRawVideosToDashVideos(videos []models.Video) ([]models.VideoForDash, error) {
	resp := make([]models.VideoForDash, 0)
	for _, video := range videos {
		dashVideo := models.VideoForDash{
			Name:          video.Name,
			Status:        video.Status,
			ReleaseDate:   video.ReleaseDate.String(),
			TokenIcon:     video.TokenIcon,
			AlbumArt:      video.AlbumArt,
			TokensOffered: video.TokensOffered,
			RevenueShared: video.RevenueShared,
		}
		resp = append(resp, dashVideo)
	}
	return resp, nil
}
