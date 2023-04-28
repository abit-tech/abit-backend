package common

const (
	CookieName                    = "token"
	RoleCookie                    = "role"
	GoogleOAuthAccessTokenRootURL = "https://oauth2.googleapis.com/token"
	GoogleOAuthFetchUserRootURL   = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token"
	HeaderKeyContentType          = "Content-Type"
	HeaderKeyAuthorization        = "Authorization"
	ContentTypeValue              = "application/x-www-form-urlencoded; charset=utf-8"

	RoleUser    = "user"
	RoleCreator = "creator"
	RoleAdmin   = "admin"

	ProviderLocal  = "local"
	ProviderGoogle = "google"

	VideoStatusPending  = "PENDING"
	VideoStatusApproved = "APPROVED"
	VideoStatusDeclined = "DECLINED"
)
