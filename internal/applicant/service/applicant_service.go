package service

import (
	// 표준 라이브러리
	"context"
	"time"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/applicant/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/applicant/repository"
)

type ApplicantService interface {
    List(ctx context.Context, page, size int, stage, keyword string) ([]model.Applicant, int, error)
    Get(ctx context.Context, id int) (*model.Applicant, error)
    UpdateStage(ctx context.Context, id int, newStage, updatedBy string) (oldStage string, updatedAt time.Time, err error)
    BulkUpdateStage(ctx context.Context, ids []int, newStage, updatedBy string) (int, error)
    GetHistory(ctx context.Context, id int) ([]model.StageHistory, error)
}

type applicantService struct {
    repo repository.ApplicantRepository
}

func NewApplicantService(repo repository.ApplicantRepository) ApplicantService {
    return &applicantService{repo: repo}
}

func (s *applicantService) List(ctx context.Context, page, size int, stage, keyword string) ([]model.Applicant, int, error) {
    offset := (page - 1) * size

    items, err := s.repo.FindAll(ctx, stage, keyword, offset, size)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.repo.CountAll(ctx, stage, keyword)

    return items, total, err
}

func (s *applicantService) Get(ctx context.Context, id int) (*model.Applicant, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *applicantService) UpdateStage(ctx context.Context, id int, newStage, updatedBy string) (string, time.Time, error) {
    a, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return "", time.Time{}, err
    }

    old := a.CurrentStage
    updatedAt, err := s.repo.UpdateStage(ctx, id, newStage, updatedBy)
    if err != nil {
        return "", time.Time{}, err
    }

    hist := &model.StageHistory{
        ApplicationID: id,
        Stage:         newStage,
        Status:        newStage,
        CreatedBy:     updatedBy,
    }
    if err := s.repo.CreateHistory(ctx, hist); err != nil {
        return "", time.Time{}, err
    }

    return old, updatedAt, nil
}

func (s *applicantService) BulkUpdateStage(ctx context.Context, ids []int, newStage, updatedBy string) (int, error) {
    cnt, err := s.repo.BulkUpdateStage(ctx, ids, newStage, updatedBy)
    if err != nil {
        return 0, err
    }

    // 히스토리도 개별 생성 (필요 시)
    for _, id := range ids {
        _ = s.repo.CreateHistory(ctx, &model.StageHistory{
            ApplicationID: id,
            Stage:         newStage,
            Status:        newStage,
            CreatedBy:     updatedBy,
        })
    }
    
    return int(cnt), nil
}

func (s *applicantService) GetHistory(ctx context.Context, id int) ([]model.StageHistory, error) {
    return s.repo.FindHistoryByApplicant(ctx, id)
}
