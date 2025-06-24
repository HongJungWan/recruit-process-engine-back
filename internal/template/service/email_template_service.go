package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/repository"
)

var ErrTemplateNotFound = errors.New("template not found")

type EmailTemplateService interface {
    List(ctx context.Context) ([]model.EmailTemplate, error)
    Get(ctx context.Context, id int) (*model.EmailTemplate, error)
    Create(ctx context.Context, name string, cfg map[string]interface{}) (*model.EmailTemplate, error)
    Update(ctx context.Context, id int, name *string, cfg *map[string]interface{}) (*model.EmailTemplate, error)
    Delete(ctx context.Context, id int) error
}

type emailTemplateService struct {
    repo repository.EmailTemplateRepository
}

func NewEmailTemplateService(repo repository.EmailTemplateRepository) EmailTemplateService {
    return &emailTemplateService{repo: repo}
}

func (s *emailTemplateService) List(ctx context.Context) ([]model.EmailTemplate, error) {
    return s.repo.FindAll(ctx)
}

func (s *emailTemplateService) Get(ctx context.Context, id int) (*model.EmailTemplate, error) {
    t, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, ErrTemplateNotFound
    }
    return t, nil
}

func (s *emailTemplateService) Create(ctx context.Context, name string, cfg map[string]interface{}) (*model.EmailTemplate, error) {
    t := &model.EmailTemplate{Name: name}
    t.Config, _ = json.Marshal(cfg)
    if err := s.repo.Create(ctx, t); err != nil {
        return nil, err
    }
    return t, nil
}

func (s *emailTemplateService) Update(ctx context.Context, id int, name *string, cfg *map[string]interface{}) (*model.EmailTemplate, error) {
    t, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, ErrTemplateNotFound
    }
    if name != nil {
        t.Name = *name
    }
    if cfg != nil {
        t.Config, _ = json.Marshal(*cfg)
    }
    if err := s.repo.Update(ctx, t); err != nil {
        return nil, err
    }
    return t, nil
}

func (s *emailTemplateService) Delete(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}
