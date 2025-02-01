package services

import (
	models "github.com/Software78/encryption-test/src/models"
	repository "github.com/Software78/encryption-test/src/repository"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repository: repo}
}

func (s *UserService) Create(user *models.User) error {
	return s.repository.Create(user)
}

func (s *UserService) Login(login *models.Login) (*models.User, error) {
	return s.repository.Login(login)
}

func (s *UserService) Register(register *models.Register) (*models.User, error) {
	return s.repository.Register(register)
}