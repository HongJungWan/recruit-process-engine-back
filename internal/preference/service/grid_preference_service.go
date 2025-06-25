package service

import (
	// 표준 라이브러리
	"context"
	"encoding/json"
	"errors"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/preference/dto/request"
	"github.com/HongJungWan/recruit-process-engine-back/internal/preference/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/preference/repository"
)

var ErrNotFound = errors.New("preference not found")

type GridPreferenceService interface {
    GetByUser(ctx context.Context, userID int) ([]model.GridPreference, error)
    Create(ctx context.Context, userID int, input request.CreateGridPreference, createdBy string) (*model.GridPreference, error)
    Update(ctx context.Context, userID, prefID int, input request.UpdateGridPreference, updatedBy string) (*model.GridPreference, error)
    Delete(ctx context.Context, userID, prefID int) error
}

type gridPreferenceService struct {
    repo repository.GridPreferenceRepository
}

func NewGridPreferenceService(repo repository.GridPreferenceRepository) GridPreferenceService {
    return &gridPreferenceService{repo: repo}
}

func (s *gridPreferenceService) GetByUser(ctx context.Context, userID int) ([]model.GridPreference, error) {
    return s.repo.FindByUser(ctx, userID)
}

func (s *gridPreferenceService) Create(ctx context.Context, userID int, input request.CreateGridPreference, createdBy string) (*model.GridPreference, error) {
    gp := &model.GridPreference{
        UserID:    userID,
        GridName:  input.GridName,
        Config:    nil,      // repo에서 Marshal
        CreatedBy: createdBy,
    }

    gp.Config, _ = json.Marshal(input.Config)
    if err := s.repo.Create(ctx, gp); err != nil {
        return nil, err
    }

    return gp, nil
}

func (s *gridPreferenceService) Update(ctx context.Context, userID, prefID int, input request.UpdateGridPreference, updatedBy string) (*model.GridPreference, error) {
    existing, err := s.repo.FindByID(ctx, prefID)
    if err != nil {
        return nil, ErrNotFound
    }
    if existing.UserID != userID {
        return nil, errors.New("forbidden")
    }

    existing.Config, _ = json.Marshal(input.Config)
    existing.UpdatedBy = &updatedBy
    if err := s.repo.Update(ctx, existing); err != nil {
        return nil, err
    }

    return existing, nil
}

func (s *gridPreferenceService) Delete(ctx context.Context, userID, prefID int) error {
    existing, err := s.repo.FindByID(ctx, prefID)
    if err != nil {
        return ErrNotFound
    }
    if existing.UserID != userID {
        return errors.New("forbidden")
    }
    
    return s.repo.Delete(ctx, prefID)
}
