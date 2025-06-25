package service

import (
	// 표준 라이브러리
	"context"

	// 내부 패키지
	repository "github.com/HongJungWan/recruit-process-engine-back/internal/user/repository"
)

type UserService interface {
    Authenticate(ctx context.Context, loginId, password string) (int, error)
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) Authenticate(ctx context.Context, loginId, password string) (int, error) {
    return s.repo.Authenticate(ctx, loginId, password)
}