package services

import (
	"blog_system/config"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"errors"

	"go.uber.org/zap"
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
		logging.Log.Warn("Language code already added database", zap.Uint("language_id", language.ID), zap.Error(err))
		return ErrAlreadyLanguage
	}

	if err := config.DB.Create(&language).Error; err != nil {
		logging.Log.Error("Failed to create language", zap.Uint("language_id", language.ID), zap.Error(err))
		return err
	}

	logging.Log.Info("Category created successfully", zap.Uint("language_id", language.ID))
	return nil
}

func (s *languageService) GetAllLanguage() ([]schemas.LanguageResponse, error) {
	var languages []models.Language
	if err := config.DB.Find(&languages).Error; err != nil {
		logging.Log.Error("Failed to get language", zap.Error(err))
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

	logging.Log.Info("Languages fetched successfully", zap.Int("count", len(response)))
	return response, nil
}

func (s *languageService) EditLanguage(ID uint, name, code string) error {
	var language models.Language
	if err := config.DB.Where("id = ?", ID).First(&language).Error; err != nil {
		logging.Log.Warn("Languages not found", zap.Uint("category_id", ID), zap.Error(err))
		return ErrLanguageNotFound
	}

	if err := config.DB.Where("code =? AND id != ?", code, ID).First(&language).Error; err == nil {
		logging.Log.Warn("Languages code already exists", zap.String("code", code))
		return ErrAlreadyCategory
	}

	if err := config.DB.Model(&language).Updates(models.Language{
		Name: name,
		Code: code,
	}).Error; err != nil {
		logging.Log.Error("Failed to update language", zap.Uint("language_id", ID), zap.Error(err))
		return err
	}

	logging.Log.Info("Languages updated successfully", zap.Uint("language_id", ID))
	return nil

}

func (s *languageService) DeleteLanguage(ID uint) error {
	var language models.Language
	if err := config.DB.Where("id = ?", ID).First(&language).Error; err != nil {
		logging.Log.Warn("Languages not found", zap.Uint("language_id", ID), zap.Error(err))
		return ErrLanguageNotFound
	}

	if err := config.DB.Unscoped().Delete(&language).Error; err != nil {
		logging.Log.Error("Failed to delete language", zap.Uint("language_id", language.ID), zap.Error(err))
		return err
	}

	logging.Log.Info("Languages deleted successfully", zap.Uint("language_id", language.ID))
	return nil
}
