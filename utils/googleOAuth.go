package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"www.github.com/abit-tech/abit-backend/common"
	"www.github.com/abit-tech/abit-backend/initializers"
)

const ()

type GoogleOauthToken struct {
	Access_token string
	Id_token     string
}

type GoogleUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
}

// GetGoogleOauthToken obtains an access token that is used to fetch Google user details
func GetGoogleOauthToken(code string) (*GoogleOauthToken, error) {
	config := initializers.AppConf

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	// code is the authorization code obtained from authorization endpoint
	values.Add("code", code)
	values.Add("client_id", config.GoogleClientID)
	values.Add("client_secret", config.GoogleClientSecret)
	values.Add("redirect_uri", config.GoogleOAuthRedirectURL)

	query := values.Encode()

	req, err := http.NewRequest("POST", common.GoogleOAuthAccessTokenRootURL, bytes.NewBufferString(query))
	if err != nil {
		// todo add log
		return nil, err
	}

	req.Header.Set(common.HeaderKeyContentType, common.ContentTypeValue)
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		// todo add log
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		// todo log error
		fmt.Printf("res from google: %v\n", res)
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		// todo log error
		return nil, err
	}

	var googleOauthTokenRes map[string]interface{}
	if err := json.Unmarshal(resBody.Bytes(), &googleOauthTokenRes); err != nil {
		// todo log error
		return nil, err
	}

	tokenBody := &GoogleOauthToken{
		Access_token: googleOauthTokenRes["access_token"].(string),
		Id_token:     googleOauthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

func GetGoogleUser(access_token string, id_token string) (*GoogleUserResult, error) {
	rootURL := fmt.Sprintf("%s=%s", common.GoogleOAuthFetchUserRootURL, access_token)

	req, err := http.NewRequest("GET", rootURL, nil)
	if err != nil {
		// todo log error
		return nil, err
	}

	req.Header.Set(common.HeaderKeyAuthorization, fmt.Sprintf("Bearer %s", id_token))
	client := http.Client{
		Timeout: time.Second & 30,
	}

	res, err := client.Do(req)
	if err != nil {
		// todo log error
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		// todo log error
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		// todo log error
		return nil, err
	}

	var googleUserRes map[string]interface{}
	if err := json.Unmarshal(resBody.Bytes(), &googleUserRes); err != nil {
		return nil, err
	}

	userBody := &GoogleUserResult{
		Id:             googleUserRes["id"].(string),
		Email:          googleUserRes["email"].(string),
		Verified_email: googleUserRes["verified_email"].(bool),
		Name:           googleUserRes["name"].(string),
		Given_name:     googleUserRes["given_name"].(string),
		Picture:        googleUserRes["picture"].(string),
		Locale:         googleUserRes["locale"].(string),
	}

	return userBody, nil
}
