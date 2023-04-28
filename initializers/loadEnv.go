package initializers

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var AppConf *Config

type Config struct {
	FrontEndOrigin string `mapstructure:"FRONTEND_ORIGIN"`

	JWTTokenSecret string        `mapstructure:"JWT_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRES_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAX_AGE"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectURL string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`

	RevenueSharingContractTemplateURI string `mapstructure:"REVENUE_SHARING_CONTRACT_TEMPLATE_URI"`
	TokenOwnershipContractTemplateURI string `mapstructure:"TOKEN_OWNERSHIP_CONTRACT_TEMPLATE_URI"`

	YoutubeSecrets YoutubeSecrets `mapstructure:"YOUTUBE_SECRETS"`

	// GoogleOAuthRedirectURLForUser    string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL_FOR_USER"`
	// GoogleOAuthRedirectURLForCreator string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL_FOR_CREATOR"`
}

type YoutubeSecrets struct {
	ApiKeyMoksh    string `mapstructure:"API_KEY_MOKSH"`
	ChannelIDMoksh string `mapstructure:"CHANNEL_ID_MOKSH"`
	ApiKeyPankaj   string `mapstructure:"API_KEY_PANKAJ"`
}

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return err
	}

	err = viper.Unmarshal(&AppConf)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return err
	}
	return nil
}
