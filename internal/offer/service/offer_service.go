package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/repository"
)

var ErrOfferNotFound = errors.New("offer not found")

type OfferService interface {
    Create(ctx context.Context, userID int, in CreateOfferInput) (*model.Offer, error)
    List(ctx context.Context, status string, page, size int) ([]model.Offer, int, error)
    GetDetail(ctx context.Context, id int) (*model.Offer, []response.ApproverStatus, error)
}

type offerService struct {
    repo repository.OfferRepository
    appRepo repository.ApprovalRepository
}

func NewOfferService(or repository.OfferRepository, ar repository.ApprovalRepository) OfferService {
    return &offerService{repo: or, appRepo: ar}
}

type CreateOfferInput struct {
    ApplicationID int
    Position      string
    Salary        int
    StartDate     time.Time
    Location      string
    Benefits      string
    LetterContent string
}

func (s *offerService) Create(ctx context.Context, userID int, in CreateOfferInput) (*model.Offer, error) {
    o := &model.Offer{
        UserID:        userID,
        ApplicationID: in.ApplicationID,
        Position:      in.Position,
        Salary:        in.Salary,
        StartDate:     in.StartDate,
        Location:      in.Location,
        Benefits:      in.Benefits,
        LetterContent: in.LetterContent,
        CreatedBy:     strconv.Itoa(userID),
    }
    if err := s.repo.Create(ctx, o); err != nil {
        return nil, err
    }
    return o, nil
}

func (s *offerService) List(ctx context.Context, status string, page, size int) ([]model.Offer, int, error) {
    offset := (page - 1) * size
    return s.repo.FindAll(ctx, status, offset, size)
}

func (s *offerService) GetDetail(ctx context.Context, id int) (*model.Offer, []response.ApproverStatus, error) {
    o, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, nil, ErrOfferNotFound
    }
    approvals, _ := s.appRepo.FindByOffer(ctx, id)
    var stats []response.ApproverStatus
    for _, ap := range approvals {
        stats = append(stats, response.ApproverStatus{
            ApproverID: ap.ApproverID,
            Status:     ap.Status,
            Comment:    ap.Comment,
        })
    }
    return o, stats, nil
}
