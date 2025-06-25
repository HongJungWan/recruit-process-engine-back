package service

import (
	// 표준 라이브러리
	"context"
	"errors"
	"time"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/repository"
)

var ErrApprovalNotFound = errors.New("approval not found")

type ApprovalService interface {
    Request(ctx context.Context, offerID int, approvers []int, requestedBy string) ([]model.Approval, error)
    ListByOffer(ctx context.Context, offerID int) ([]model.Approval, error)
    Process(ctx context.Context, id int, status, comment, decidedBy string) (time.Time, error)
}

type approvalService struct {
    repo repository.ApprovalRepository
}

func NewApprovalService(ar repository.ApprovalRepository) ApprovalService {
    return &approvalService{repo: ar}
}

func (s *approvalService) Request(ctx context.Context, offerID int, approvers []int, requestedBy string) ([]model.Approval, error) {
    return s.repo.CreateBulk(ctx, offerID, approvers, requestedBy)
}

func (s *approvalService) ListByOffer(ctx context.Context, offerID int) ([]model.Approval, error) {
    return s.repo.FindByOffer(ctx, offerID)
}

func (s *approvalService) Process(ctx context.Context, id int, status, comment, decidedBy string) (time.Time, error) {
    _, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return time.Time{}, ErrApprovalNotFound
    }

    decidedAt, err := s.repo.UpdateStatus(ctx, id, status, comment, decidedBy)

    return decidedAt, err
}
