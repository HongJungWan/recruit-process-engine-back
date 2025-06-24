package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/model"
)

type OfferRepository interface {
    Create(ctx context.Context, o *model.Offer) error
    FindAll(ctx context.Context, status string, offset, limit int) ([]model.Offer, int, error)
    FindByID(ctx context.Context, id int) (*model.Offer, error)
}

type offerRepo struct {
    db    *sqlx.DB
    sb    sq.StatementBuilderType
    table string
}

func NewOfferRepository(db *sqlx.DB) OfferRepository {
    return &offerRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: (&model.Offer{}).TableName(),
    }
}

func (r *offerRepo) Create(ctx context.Context, o *model.Offer) error {
    qb := r.sb.Insert(r.table).
        Columns("user_id","application_id","position","salary","start_date","location","benefits","letter_content","status","created_by").
        Values(o.UserID,o.ApplicationID,o.Position,o.Salary,o.StartDate,o.Location,o.Benefits,o.LetterContent,"PENDING",o.CreatedBy).
        Suffix("RETURNING offer_id,status,created_at")
    sqlStr, args, _ := qb.ToSql()
    return r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(o)
}

func (r *offerRepo) FindAll(ctx context.Context, status string, offset, limit int) ([]model.Offer, int, error) {
    qb := r.sb.Select("*").From(r.table)
    countQ := r.sb.Select("COUNT(*)").From(r.table)
    if status != "" {
        qb = qb.Where(sq.Eq{"status": status})
        countQ = countQ.Where(sq.Eq{"status": status})
    }
    qb = qb.Limit(uint64(limit)).Offset(uint64(offset))
    // fetch items
    sqlStr, args, _ := qb.ToSql()
    var list []model.Offer
    if err := r.db.SelectContext(ctx, &list, sqlStr, args...); err != nil {
        return nil, 0, err
    }
    // count total
    cntStr, cntArgs, _ := countQ.ToSql()
    var total int
    if err := r.db.GetContext(ctx, &total, cntStr, cntArgs...); err != nil {
        return nil, 0, err
    }
    return list, total, nil
}

func (r *offerRepo) FindByID(ctx context.Context, id int) (*model.Offer, error) {
    qb := r.sb.Select("*").From(r.table).Where(sq.Eq{"offer_id": id})
    sqlStr, args, _ := qb.ToSql()
    var o model.Offer
    if err := r.db.GetContext(ctx, &o, sqlStr, args...); err != nil {
        return nil, err
    }
    return &o, nil
}
