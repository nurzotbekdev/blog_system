package schemas

type LikeRequest struct {
	VideoID uint `json:"video_id"`
	IsLike  bool `json:"is_like"`
}

type LikeEdite struct {
	IsLike bool `json:"is_like"`
}

type UserLikeResponse struct {
	ID      uint   `json:"id"`
	VideoID uint   `json:"video_id"`
	Title   string `json:"title"`
	IsLike  bool   `json:"is_like"`
}
