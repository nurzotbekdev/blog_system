package validators

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func isValidImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	return allowed[ext]
}

func ValidateChannel(name string, profile *multipart.FileHeader, banner *multipart.FileHeader) error {

	if name == "" {
		return errors.New("name is required")
	}

	if len(name) < 3 {
		return errors.New("name too short")
	}

	if profile != nil {
		if profile.Size > 2*1024*1024 {
			return errors.New("profile image too large (max 2MB)")
		}

		if !isValidImage(profile.Filename) {
			return errors.New("invalid profile image format (only jpg, png allowed)")
		}
	}

	if banner != nil {
		if banner.Size > 5*1024*1024 {
			return errors.New("banner image too large (max 5MB)")
		}

		if !isValidImage(banner.Filename) {
			return errors.New("invalid banner image format (only jpg, png allowed)")
		}
	}

	return nil
}
