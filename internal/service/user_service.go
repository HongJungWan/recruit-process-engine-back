package service

import (
    "context"
    "errors"
    "github.com/HongJungWan/recruit-process-engine-back/internal/models"
    "github.com/HongJungWan/recruit-process-engine-back/internal/repository"
)

type UserService interface {
    Register(ctx context.Context, email, password, name string) (int, error)
    GetByID(ctx context.Context, id int) (*models.User, error)
    Login(ctx context.Context, email, password string) (*models.User, error)
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, email, password, name string) (int, error) {
    existUser, _ := s.repo.GetByEmail(ctx, email)
    if existUser != nil {
        return 0, errors.New("email already in use")
    }

    hashedPassword := password

    u := &models.User{
        Email:    email,
        Password: hashedPassword,
        Name:     name,
    }

    newID, err := s.repo.Create(ctx, u)
    if err != nil {
        return 0, err
    }
    return newID, nil
}

func (s *userService) GetByID(ctx context.Context, id int) (*models.User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *userService) Login(ctx context.Context, email, password string) (*models.User, error) {
    u, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if u == nil {
        return nil, errors.New("user not found")
    }
    if u.Password != password {
        return nil, errors.New("wrong credentials")
    }
    return u, nil
}
