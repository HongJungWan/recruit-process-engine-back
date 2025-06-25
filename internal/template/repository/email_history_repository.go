package repository

import (
	// 표준 라이브러리
	"context"

	// 서드파티(외부) 라이브러리
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
)

type EmailHistoryRepository interface {
    Create(ctx context.Context, h *model.EmailHistory) error
    FindAll(ctx context.Context, applicantID, offerID *int, offset, limit int) ([]model.EmailHistory, error)
    FindByID(ctx context.Context, id int) (*model.EmailHistory, error)
}

type emailHistoryRepo struct {
    db        *sqlx.DB
    sb        sq.StatementBuilderType
    table     string
}

func NewEmailHistoryRepository(db *sqlx.DB) EmailHistoryRepository {
    return &emailHistoryRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: (&model.EmailHistory{}).TableName(),
    }
}

func (r *emailHistoryRepo) Create(ctx context.Context, h *model.EmailHistory) error {
    qb := r.sb.
        Insert(r.table).
            Columns(
                "user_id",
                "application_id",
                "offer_id",
                "template_id",
                "title",
                "body",
                "created_by").
            Values(
                h.UserID, 
                h.ApplicationID, 
                h.OfferID, 
                h.TemplateID, 
                h.Title, 
                h.Body, 
                h.CreatedBy).
        Suffix("RETURNING email_id, created_at")
    
    sqlStr, args, _ := qb.ToSql()
    
    return r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(h)
}

func (r *emailHistoryRepo) FindAll(ctx context.Context, applicantID, offerID *int, offset, limit int) ([]model.EmailHistory, error) {
    qb := r.sb.
        Select(
            "email_id",
            "title",
            "created_at").
        From(r.table)

    if applicantID != nil {
        qb = qb.Where(sq.Eq{"application_id": *applicantID})
    }
    if offerID != nil {
        qb = qb.Where(sq.Eq{"offer_id": *offerID})
    }

    qb = qb.Limit(uint64(limit)).Offset(uint64(offset))
    
    var list []model.EmailHistory
    sqlStr, args, _ := qb.ToSql()
    if err := r.db.SelectContext(ctx, &list, sqlStr, args...); err != nil {
        return nil, err
    }

    return list, nil
}

func (r *emailHistoryRepo) FindByID(ctx context.Context, id int) (*model.EmailHistory, error) {
    sqlStr, args, _ := r.sb.
        Select("*").
        From(r.table).
        Where(sq.Eq{"email_id": id}).
        ToSql()
    
    var h model.EmailHistory
    if err := r.db.GetContext(ctx, &h, sqlStr, args...); err != nil {
        return nil, err
    }
    
    return &h, nil
}
