package models

type YoutubeVideo struct {
	Views        int
	AdRevenue    float32
	GrossRevenue float32
	CPM          float32
}

type YoutubeChannel struct {
	Name            string
	TotalViews      int
	TotalLikes      int
	SubscriberCount int
	VideoCount      int
}
