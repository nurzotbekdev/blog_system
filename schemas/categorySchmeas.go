package schemas

import "time"

type CategorySchemas struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
