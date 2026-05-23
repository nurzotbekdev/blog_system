package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"
)

type LanguageService interface {
	CreateLanguage(language models.Language) error
	GetAllLanguage() ([]schemas.LanguageResponse, error)
	EditLanguage(ID uint, name, code string) error
	DeleteLanguage(ID uint) error
}

type languageService struct{}

func NewLanguageServices() LanguageService {
	return &languageService{}
}

var (
	ErrAlreadyLanguage  = errors.New("Language already code added database")
	ErrLanguageNotFound = errors.New("Language not found")
)

func (s *languageService) CreateLanguage(language models.Language) error {
	var existing models.Language
	if err := config.DB.Where("code = ?", language.Code).First(&existing).Error; err == nil {
		return ErrAlreadyLanguage
	}

	if err := config.DB.Create(&language).Error; err != nil {
		return err
	}

	return nil
}

func (s *languageService) GetAllLanguage() ([]schemas.LanguageResponse, error) {
	var languages []models.Language
	if err := config.DB.Find(&languages).Error; err != nil {
		return nil, err
	}

	var response []schemas.LanguageResponse
	for _, l := range languages {
		response = append(response, schemas.LanguageResponse{
			ID:        l.ID,
			Name:      l.Name,
			Code:      l.Code,
			CreatedAt: l.CreatedAt,
		})
	}

	return response, nil
}

func (s *languageService) EditLanguage(ID uint, name, code string) error {
	var language models.Language
	if err := config.DB.Where("id = ?", ID).First(&language).Error; err != nil {
		return ErrLanguageNotFound
	}

	if err := config.DB.Where("code =? AND id != ?", code, ID).First(&language).Error; err == nil {
		return ErrAlreadyCategory
	}

	if err := config.DB.Model(&language).Updates(models.Language{
		Name: name,
		Code: code,
	}).Error; err != nil {
		return err
	}

	return nil

}

func (s *languageService) DeleteLanguage(ID uint) error {
	var language models.Language
	if err := config.DB.Where("id = ?", ID).First(&language).Error; err != nil {
		return ErrLanguageNotFound
	}

	if err := config.DB.Unscoped().Delete(&language).Error; err != nil {
		return err
	}

	return nil
}
