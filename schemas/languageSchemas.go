package schemas

import "time"

type LanguageSchemas struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type LanguageResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}
