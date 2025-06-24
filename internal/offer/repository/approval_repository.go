package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/model"
)

type ApprovalRepository interface {
    CreateBulk(ctx context.Context, offerID int, approvers []int, requestedBy string) ([]model.Approval, error)
    FindByOffer(ctx context.Context, offerID int) ([]model.Approval, error)
    UpdateStatus(ctx context.Context, id int, status, comment, decidedBy string) (time.Time, error)
    FindByID(ctx context.Context, id int) (*model.Approval, error)
}

type approvalRepo struct {
    db    *sqlx.DB
    sb    sq.StatementBuilderType
    table string
}

func NewApprovalRepository(db *sqlx.DB) ApprovalRepository {
    return &approvalRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: (&model.Approval{}).TableName(),
    }
}

func (r *approvalRepo) CreateBulk(ctx context.Context, offerID int, approvers []int, requestedBy string) ([]model.Approval, error) {
    var created []model.Approval
    for _, a := range approvers {
        qb := r.sb.Insert(r.table).
            Columns("offer_id","approver_id","created_by").
            Values(offerID, a, requestedBy).
            Suffix("RETURNING approval_id,offer_id,approver_id,status,requested_at")
        sqlStr, args, _ := qb.ToSql()
        var ap model.Approval
        if err := r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(&ap); err != nil {
            return nil, err
        }
        created = append(created, ap)
    }
    return created, nil
}

func (r *approvalRepo) FindByOffer(ctx context.Context, offerID int) ([]model.Approval, error) {
    qb := r.sb.Select("*").From(r.table).Where(sq.Eq{"offer_id": offerID})
    sqlStr, args, _ := qb.ToSql()
    var list []model.Approval
    if err := r.db.SelectContext(ctx, &list, sqlStr, args...); err != nil {
        return nil, err
    }
    return list, nil
}

func (r *approvalRepo) UpdateStatus(ctx context.Context, id int, status, comment, decidedBy string) (time.Time, error) {
    qb := r.sb.Update(r.table).
        Set("status", status).
        Set("comment", comment).
        Set("decided_at", sq.Expr("now()")).
        Where(sq.Eq{"approval_id": id}).
        Suffix("RETURNING decided_at")
    sqlStr, args, _ := qb.ToSql()
    var decidedAt time.Time
    if err := r.db.GetContext(ctx, &decidedAt, sqlStr, args...); err != nil {
        return time.Time{}, err
    }
    return decidedAt, nil
}

func (r *approvalRepo) FindByID(ctx context.Context, id int) (*model.Approval, error) {
    qb := r.sb.Select("*").From(r.table).Where(sq.Eq{"approval_id": id})
    sqlStr, args, _ := qb.ToSql()
    var ap model.Approval
    if err := r.db.GetContext(ctx, &ap, sqlStr, args...); err != nil {
        return nil, err
    }
    return &ap, nil
}
