package service

import (
	"context"
	"errors"
	"time"

	sr "github.com/HongJungWan/recruit-process-engine-back/internal/session/repository"
	um "github.com/HongJungWan/recruit-process-engine-back/internal/user/model"
	ur "github.com/HongJungWan/recruit-process-engine-back/internal/user/repository"
)

type SessionService interface {
    Login(ctx context.Context, loginID, password string) (user *um.User, token string, err error)
    Logout(ctx context.Context, token string) error
}

type sessionService struct {
    userRepo    ur.UserRepository
    sessionRepo sr.SessionRepository
    ttl         time.Duration
}

func NewSessionService(u ur.UserRepository, s sr.SessionRepository, ttl time.Duration) SessionService {
    return &sessionService{userRepo: u, sessionRepo: s, ttl: ttl}
}

// feat: 로그인
func (s *sessionService) Login(ctx context.Context, loginID, password string) (*um.User, string, error) {
    user, err := s.userRepo.GetByEmail(ctx, loginID)
    if err != nil || user == nil || user.Password != password {
        return nil, "", errors.New("invalid credentials")
    }
    token, err := s.sessionRepo.Create(ctx, user.UserID, s.ttl)
    return user, token, err
}

// feat: 로그아웃
func (s *sessionService) Logout(ctx context.Context, token string) error {
    return s.sessionRepo.Delete(ctx, token)
}
