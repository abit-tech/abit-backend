package initializers

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleConfig *oauth2.Config

func SetupGoogleOauth() {
	clientID := AppConf.GoogleClientID
	clientSecret := AppConf.GoogleClientSecret
	googleRedirectURL := AppConf.GoogleOAuthRedirectURL

	if clientID == "" || clientSecret == "" || googleRedirectURL == "" {
		panic("could not find google credentials")
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  googleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	GoogleConfig = conf
}
