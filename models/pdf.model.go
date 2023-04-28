package models

type RevenueSharingPdfData struct {
	VideoID     string `json:"videoID"`
	CreatorID   string `json:"creatorID"`
	CreatorName string `json:"creatorName"`
	VideoName   string `json:"videoName"`
	ReleaseDate string `json:"releaseDate"`
	// RevenueShared must be per token and not total
	RevenueShared  string `json:"revenueShared"`
	TokensReleased string `json:"tokensReleased"`
	TokenPrice     string `json:"tokenPrice"`
}

type OwnershipPdfData struct {
	VideoID      string `json:"videoID"`
	TokenID      string `json:"tokenID"`
	OwnerID      string `json:"ownerID"`
	CreatorName  string `json:"creatorName"`
	OwnerName    string `json:"ownerName"`
	VideoName    string `json:"videoName"`
	ReleaseDate  string `json:"releaseDate"`
	RevenueShare string `json:"revenueShare"`
	TokenPrice   string `json:"tokenPrice"`
}
