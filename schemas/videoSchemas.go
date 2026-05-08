package schemas

import "mime/multipart"

type CreateVideoRequest struct {
	LanguageID  uint   `form:"language_id"`
	CategoryID  uint   `form:"category_id"`
	Title       string `form:"title"`
	Description string `form:"description"`

	FilePath      *multipart.FileHeader `form:"file_path"`
	ThumbnailPath *multipart.FileHeader `form:"thumbnail_path"`
}
