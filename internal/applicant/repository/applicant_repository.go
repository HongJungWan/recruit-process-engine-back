package repository

import (
	"context"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/applicant/model"
)

type ApplicantRepository interface {
    FindAll(ctx context.Context, stage, keyword string, offset, limit int) ([]model.Applicant, error)
    CountAll(ctx context.Context, stage, keyword string) (int, error)
    FindByID(ctx context.Context, id int) (*model.Applicant, error)
    UpdateStage(ctx context.Context, id int, newStage, updatedBy string) (time.Time, error)
    BulkUpdateStage(ctx context.Context, ids []int, newStage, updatedBy string) (int64, error)
    CreateHistory(ctx context.Context, h *model.StageHistory) error
    FindHistoryByApplicant(ctx context.Context, id int) ([]model.StageHistory, error)
}

type applicantRepo struct {
    db        *sqlx.DB
    sb        sq.StatementBuilderType
    table     string
    histTable string
}

func NewApplicantRepository(db *sqlx.DB) ApplicantRepository {
    return &applicantRepo{
        db:        db,
        sb:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table:     (&model.Applicant{}).TableName(),
        histTable: (&model.StageHistory{}).TableName(),
    }
}

func (r *applicantRepo) FindAll(ctx context.Context, stage, keyword string, offset, limit int) ([]model.Applicant, error) {
    qb := r.sb.
        Select("application_id","name","email","current_stage").
        From(r.table)
    if stage != "" {
        qb = qb.Where(sq.Eq{"current_stage": stage})
    }
    if keyword != "" {
        kw := "%" + strings.ToLower(keyword) + "%"
        qb = qb.Where(sq.Or{
            sq.Like{"LOWER(name)":  kw},
            sq.Like{"LOWER(email)": kw},
        })
    }
    qb = qb.Limit(uint64(limit)).Offset(uint64(offset))

    sqlStr, args, _ := qb.ToSql()
    var list []model.Applicant
    if err := r.db.SelectContext(ctx, &list, sqlStr, args...); err != nil {
        return nil, err
    }
    return list, nil
}

func (r *applicantRepo) CountAll(ctx context.Context, stage, keyword string) (int, error) {
    qb := r.sb.Select("COUNT(*)").From(r.table)
    if stage != "" {
        qb = qb.Where(sq.Eq{"current_stage": stage})
    }
    if keyword != "" {
        kw := "%" + strings.ToLower(keyword) + "%"
        qb = qb.Where(sq.Or{
            sq.Like{"LOWER(name)":  kw},
            sq.Like{"LOWER(email)": kw},
        })
    }
    sqlStr, args, _ := qb.ToSql()
    var total int
    if err := r.db.GetContext(ctx, &total, sqlStr, args...); err != nil {
        return 0, err
    }
    return total, nil
}

func (r *applicantRepo) FindByID(ctx context.Context, id int) (*model.Applicant, error) {
    qb := r.sb.Select("*").From(r.table).Where(sq.Eq{"application_id": id})
    sqlStr, args, _ := qb.ToSql()
    var a model.Applicant
    if err := r.db.GetContext(ctx, &a, sqlStr, args...); err != nil {
        return nil, err
    }
    return &a, nil
}

func (r *applicantRepo) UpdateStage(ctx context.Context, id int, newStage, updatedBy string) (time.Time, error) {
    qb := r.sb.Update(r.table).
        Set("current_stage", newStage).
        Set("updated_by", updatedBy).
        Set("updated_at", sq.Expr("now()")).
        Where(sq.Eq{"application_id": id}).
        Suffix("RETURNING updated_at")
    sqlStr, args, _ := qb.ToSql()

    var updatedAt time.Time
    if err := r.db.GetContext(ctx, &updatedAt, sqlStr, args...); err != nil {
        return time.Time{}, err
    }
    return updatedAt, nil
}

func (r *applicantRepo) BulkUpdateStage(ctx context.Context, ids []int, newStage, updatedBy string) (int64, error) {
    qb := r.sb.Update(r.table).
        Set("current_stage", newStage).
        Set("updated_by", updatedBy).
        Set("updated_at", sq.Expr("now()")).
        Where(sq.Eq{"application_id": ids})
    sqlStr, args, _ := qb.ToSql()
    res, err := r.db.ExecContext(ctx, sqlStr, args...)
    if err != nil {
        return 0, err
    }
    return res.RowsAffected()
}

func (r *applicantRepo) CreateHistory(ctx context.Context, h *model.StageHistory) error {
    qb := r.sb.Insert(r.histTable).
        Columns("application_id","stage","status","created_by").
        Values(h.ApplicationID, h.Stage, h.Status, h.CreatedBy).
        Suffix("RETURNING history_id, created_at")
    sqlStr, args, _ := qb.ToSql()

    return r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(h)
}

func (r *applicantRepo) FindHistoryByApplicant(ctx context.Context, id int) ([]model.StageHistory, error) {
    qb := r.sb.Select("*").
        From(r.histTable).
        Where(sq.Eq{"application_id": id}).
        OrderBy("created_at")
    sqlStr, args, _ := qb.ToSql()

    var hs []model.StageHistory
    if err := r.db.SelectContext(ctx, &hs, sqlStr, args...); err != nil {
        return nil, err
    }
    return hs, nil
}
