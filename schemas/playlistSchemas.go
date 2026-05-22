package schemas

import "time"

type PlaylistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

type PlaylistResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	Thumbnail   string    `json:"thumbnail"`
	VideoCount  int64     `json:"video_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type PlaylistUpdate struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPrivate   *bool   `json:"is_private"`
}
