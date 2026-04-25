package services

import (
	"blog_system/config"
	"blog_system/models"
)

type UserService interface {
	SignIn(user models.User) (models.User, error)
}

type userService struct{}

func NewUserServices() UserService {
	return &userService{}
}

func (s *userService) SignIn(user models.User) (models.User, error) {
	var existingUser models.User

	err := config.DB.Where(models.User{GoogleID: user.GoogleID}).
		Assign(models.User{
			FullName:     user.FullName,
			ProfileImage: user.ProfileImage,
			Email:        user.Email,
		}).
		FirstOrCreate(&existingUser).Error

	return existingUser, err
}
