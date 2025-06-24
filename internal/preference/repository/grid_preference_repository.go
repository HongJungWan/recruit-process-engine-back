package repository

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/preference/model"
)

type GridPreferenceRepository interface {
    FindByUser(ctx context.Context, userID int) ([]model.GridPreference, error)
    FindByID(ctx context.Context, id int) (*model.GridPreference, error)
    Create(ctx context.Context, gp *model.GridPreference) error
    Update(ctx context.Context, gp *model.GridPreference) error
    Delete(ctx context.Context, id int) error
}

type gridPreferenceRepo struct {
    db    *sqlx.DB
    sb    sq.StatementBuilderType
    table string
}

func NewGridPreferenceRepository(db *sqlx.DB) GridPreferenceRepository {
    return &gridPreferenceRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: (&model.GridPreference{}).TableName(),
    }
}

func (r *gridPreferenceRepo) FindByUser(ctx context.Context, userID int) ([]model.GridPreference, error) {
    qb := r.sb.
        Select("*").
        From(r.table).
        Where(sq.Eq{"user_id": userID})
    sqlStr, args, _ := qb.ToSql()

    var list []model.GridPreference
    if err := r.db.SelectContext(ctx, &list, sqlStr, args...); err != nil {
        return nil, err
    }
    return list, nil
}

func (r *gridPreferenceRepo) FindByID(ctx context.Context, id int) (*model.GridPreference, error) {
    qb := r.sb.
        Select("*").
        From(r.table).
        Where(sq.Eq{"preference_id": id})
    sqlStr, args, _ := qb.ToSql()

    var gp model.GridPreference
    if err := r.db.GetContext(ctx, &gp, sqlStr, args...); err != nil {
        return nil, err
    }
    return &gp, nil
}

func (r *gridPreferenceRepo) Create(ctx context.Context, gp *model.GridPreference) error {
    data, _ := json.Marshal(gp.Config)
    qb := r.sb.
        Insert(r.table).
        Columns("user_id","grid_name","config","created_by").
        Values(gp.UserID, gp.GridName, data, gp.CreatedBy).
        Suffix("RETURNING preference_id, created_at")
    sqlStr, args, _ := qb.ToSql()

    return r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(gp)
}

func (r *gridPreferenceRepo) Update(ctx context.Context, gp *model.GridPreference) error {
    data, _ := json.Marshal(gp.Config)
    qb := r.sb.
        Update(r.table).
        Set("config", data).
        Set("updated_by", gp.UpdatedBy).
        Set("updated_at", sq.Expr("now()")).
        Where(sq.Eq{"preference_id": gp.PreferenceID}).
        Suffix("RETURNING updated_at")
    sqlStr, args, _ := qb.ToSql()

    return r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(&gp.UpdatedAt)
}

func (r *gridPreferenceRepo) Delete(ctx context.Context, id int) error {
    qb := r.sb.
        Delete(r.table).
        Where(sq.Eq{"preference_id": id})
    sqlStr, args, _ := qb.ToSql()
    _, err := r.db.ExecContext(ctx, sqlStr, args...)
    return err
}
