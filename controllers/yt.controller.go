package controllers

import (
	"errors"
	"fmt"

	"google.golang.org/api/youtube/v3"
	"www.github.com/abit-tech/abit-backend/initializers"
)

// channel list api returns the basic metadata about channel

func GetChannelDataFromYoutube(yt *youtube.Service, channelID string) error {
	conf := initializers.AppConf.YoutubeSecrets

	call := yt.Channels.List([]string{"snippet, contentDetails"}).Id(conf.ChannelIDMoksh)

	resp, err := call.Do()
	if err != nil {
		return err
	}

	if resp == nil || resp.Items == nil {
		fmt.Println("nil response received")
	}

	// rawData := resp.Items[0]
	// channelData := models.YoutubeChannel{
	// 	TotalViews:      int(rawData.Statistics.ViewCount),
	// 	SubscriberCount: int(rawData.Statistics.SubscriberCount),
	// }
	return nil
}

// Video list api returns the view and likes for a video

func GetVideoDataFromYoutube(yt *youtube.Service, channelID string, videoID string) ([]*youtube.SearchResult, error) {
	call := yt.Search.List([]string{"snippet"}).
		ChannelId(channelID).Type("video").
		MaxResults(20).Order("date")

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("nil response received")
	}

	// we could extract just the video ID from the result and return a clean string, but in doing so
	// we also lose out on a bunch of extra information (such as the video title, published at, etc)
	// which might later come in handy when we try to enrich the FE with more details about the clip
	var videoIDs []*youtube.SearchResult
	for _, item := range resp.Items {
		if item.Id.Kind == "youtube#video" {
			videoIDs = append(videoIDs, item)
		}
	}

	return videoIDs, nil
}
