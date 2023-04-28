package initializers

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var YoutubeClient *youtube.Service

func SetupYoutubeClient() {
	key := AppConf.YoutubeSecrets.ApiKeyMoksh
	yt, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		log.Fatal("failed to instantiate youtube client")
	}
	YoutubeClient = yt
}
