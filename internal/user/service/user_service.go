package service

import (
	"context"

	repository "github.com/HongJungWan/recruit-process-engine-back/internal/user/repository"
)

type UserService interface {
    Authenticate(ctx context.Context, username, password string) (int, error)
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (int, error) {
    return s.repo.Authenticate(ctx, username, password)
}