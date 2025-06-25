package service

import (
	// 표준 라이브러리
	"context"
	"time"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/repository"
)

type EmailHistoryService interface {
    Send(ctx context.Context, h *model.EmailHistory) (*model.EmailHistory, error)
    List(ctx context.Context, applicantID, offerID *int, page, size int) ([]model.EmailHistory, error)
    Get(ctx context.Context, id int) (*model.EmailHistory, error)
}

type emailHistoryService struct {
    repo repository.EmailHistoryRepository
}

func NewEmailHistoryService(repo repository.EmailHistoryRepository) EmailHistoryService {
    return &emailHistoryService{repo: repo}
}

func (s *emailHistoryService) Send(ctx context.Context, h *model.EmailHistory) (*model.EmailHistory, error) {
    h.CreatedAt = time.Time{}
    if err := s.repo.Create(ctx, h); err != nil {
        return nil, err
    }

    return h, nil
}

func (s *emailHistoryService) List(ctx context.Context, applicantID, offerID *int, page, size int) ([]model.EmailHistory, error) {
    offset := (page - 1) * size

    return s.repo.FindAll(ctx, applicantID, offerID, offset, size)
}

func (s *emailHistoryService) Get(ctx context.Context, id int) (*model.EmailHistory, error) {
    return s.repo.FindByID(ctx, id)
}
