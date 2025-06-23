package service

import (
	repository "github.com/HongJungWan/recruit-process-engine-back/internal/user/repository"
)

type UserService interface {
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}
