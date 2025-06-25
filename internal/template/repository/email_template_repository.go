package repository

import (
	// 표준 라이브러리
	"context"
	"encoding/json"

	// 서드파티(외부) 라이브러리
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
)

type EmailTemplateRepository interface {
    FindAll(ctx context.Context) ([]model.EmailTemplate, error)
    FindByID(ctx context.Context, id int) (*model.EmailTemplate, error)
    Create(ctx context.Context, t *model.EmailTemplate) error
    Update(ctx context.Context, t *model.EmailTemplate) error
    Delete(ctx context.Context, id int) error
}

type emailTemplateRepo struct {
    db    *sqlx.DB
    sb    sq.StatementBuilderType
    table string
}

func NewEmailTemplateRepository(db *sqlx.DB) EmailTemplateRepository {
    return &emailTemplateRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: (&model.EmailTemplate{}).TableName(),
    }
}

func (r *emailTemplateRepo) FindAll(ctx context.Context) ([]model.EmailTemplate, error) {
    sqlStr, args, _ := r.sb.
    Select("*").
    From(r.table).
    ToSql()

    var list []model.EmailTemplate
    if err := r.db.SelectContext(ctx, &list, sqlStr, args...); err != nil {
        return nil, err
    }

    return list, nil
}

func (r *emailTemplateRepo) FindByID(ctx context.Context, id int) (*model.EmailTemplate, error) {
    sqlStr, args, _ := r.sb.
    Select("*").
    From(r.table).
    Where(sq.Eq{"id": id}).
    ToSql()

    var t model.EmailTemplate
    if err := r.db.GetContext(ctx, &t, sqlStr, args...); err != nil {
        return nil, err
    }

    return &t, nil
}

func (r *emailTemplateRepo) Create(ctx context.Context, t *model.EmailTemplate) error {
    data, _ := json.Marshal(t.Config)

    sqlStr, args, _ := r.sb.
    Insert(r.table).
        Columns("name", "config").
        Values(t.Name, data).
    Suffix("RETURNING id, created_at").
    ToSql()

    return r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(t)
}

func (r *emailTemplateRepo) Update(ctx context.Context, t *model.EmailTemplate) error {
    sqlStr, args, _ := r.sb.
    Update(r.table).
        Set("name", t.Name).
        Set("config", t.Config).
    Where(sq.Eq{"id": t.ID}).
    Suffix("RETURNING created_at").
    ToSql()

    return r.db.QueryRowxContext(ctx, sqlStr, args...).StructScan(&t.CreatedAt)
}

func (r *emailTemplateRepo) Delete(ctx context.Context, id int) error {
    sqlStr, args, _ := r.sb.
    Delete(r.table).
    Where(sq.Eq{"id": id}).
    ToSql()

    _, err := r.db.ExecContext(ctx, sqlStr, args...)
    
    return err
}
