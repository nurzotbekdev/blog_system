package schemas

type SubscriptionSchemas struct {
	ChannelID uint `json:"channel_id"`
}

type SubscriptionResponse struct {
	ID               uint   `json:"id"`
	ChannelID        uint   `json:"channel_id"`
	Name             string `json:"name"`
	ProfileImage     string `json:"profile_image"`
	BannerImage      string `json:"banner_image"`
	TotalSubscribers uint   `json:"total_subscribers"`
	TotalVideos      uint   `json:"total_videos"`
	TotalViews       uint   `json:"total_views"`
}

type ChannelSubscriptionResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	FullName     string `json:"full_name"`
	ProfileImage string `json:"profile_image"`
}

type SubscriptionStatus struct {
	Total     int64 `json:"total"`
	Today     int64 `json:"today"`
	ThisWeek  int64 `json:"this_week"`
	ThisMonth int64 `json:"this_month"`
}
