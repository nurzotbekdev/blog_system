package schemas

import "time"

type CommentSchemas struct {
	VideoID  uint   `json:"video_id"`
	Content  string `json:"content"`
	ParentID *uint  `json:"parent_id"`
}

type CommentResponse struct {
	ID        uint              `json:"id"`
	UserID    uint              `json:"user_id"`
	FullName  string            `json:"full_name"`
	Avatar    string            `json:"avatar"`
	Content   string            `json:"content"`
	LikeCount int64             `json:"like_count"`
	CreatedAt time.Time         `json:"created_at"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}

type CommentListResponse struct {
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
	Data       []CommentResponse `json:"data"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type CommentStatsResponse struct {
	CommentID    uint  `json:"comment_id"`
	LikeCount    int64 `json:"like_count"`
	DislikeCount int64 `json:"dislike_count"`
}
