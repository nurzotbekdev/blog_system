package schemas

import "time"

type CommentReactionRequest struct {
	IsLike bool `json:"is_like"`
}

type CommentReactionResponse struct {
	ID        uint      `json:"id"`
	IsLike    bool      `json:"is_like"`
	CreatedAt time.Time `json:"created_at"`
}
