package helper

import (
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SaveFile(file *multipart.FileHeader, folder string) (string, error) {
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)

	filename := uuid.New().String() + ext

	uploadPath := filepath.Join(folder, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(uploadPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", err
	}

	return uploadPath, nil
}

func RemoveFile(paths ...string) {
	for _, p := range paths {
		if p != "" {
			os.Remove(p)
		}
	}
}
