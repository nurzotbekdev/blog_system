package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"
)

type CategoryService interface {
	CreateCategory(category models.Category) error
	GetAllCategory() ([]schemas.CategoryResponse, error)
	EditCategory(ID uint, name string) error
	DeleteCategory(ID uint) error
}

type categoryService struct{}

func NewCategoryServices() CategoryService {
	return &categoryService{}
}

var (
	ErrAlreadyCategory  = errors.New("Category already added database")
	ErrCategoryNotFound = errors.New("Category not found")
)

func (s *categoryService) CreateCategory(category models.Category) error {
	var existing models.Category

	if err := config.DB.
		Where("name = ?", category.Name).
		First(&existing).Error; err == nil {

		return ErrAlreadyCategory
	}

	if err := config.DB.Create(&category).Error; err != nil {
		return err
	}

	return nil
}

func (s *categoryService) GetAllCategory() ([]schemas.CategoryResponse, error) {
	var categorys []models.Category
	if err := config.DB.Find(&categorys).Error; err != nil {
		return nil, err
	}

	var response []schemas.CategoryResponse
	for _, c := range categorys {
		response = append(response, schemas.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
		})
	}

	return response, nil
}

func (s *categoryService) EditCategory(ID uint, name string) error {
	var category models.Category

	if err := config.DB.Where("id = ?", ID).First(&category).Error; err != nil {
		return ErrCategoryNotFound
	}

	var existing models.Category
	if err := config.DB.
		Where("name = ? AND id != ?", name, ID).
		First(&existing).Error; err == nil {

		return ErrAlreadyCategory
	}

	if err := config.DB.Model(&category).Update("name", name).Error; err != nil {

		return err
	}

	return nil
}

func (s *categoryService) DeleteCategory(ID uint) error {
	var category models.Category
	if err := config.DB.Where("id = ?", ID).First(&category).Error; err != nil {

		return ErrCategoryNotFound
	}

	if err := config.DB.Unscoped().Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
