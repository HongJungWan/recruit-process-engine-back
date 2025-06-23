package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session/model"
)

type SessionRepository interface {
    Create(ctx context.Context, userID int, ttl time.Duration) (string, error)
    Delete(ctx context.Context, token string) error
    GetByToken(ctx context.Context, token string) (*model.Session, error)
}

type sessionRepo struct {
    db    *sqlx.DB
    sb    sq.StatementBuilderType
    table string
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
    return &sessionRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: (&model.Session{}).TableName(),
    }
}

// 새 세션 생성
func (r *sessionRepo) Create(ctx context.Context, userID int, ttl time.Duration) (string, error) {
    token := uuid.New().String()
    expires := time.Now().Add(ttl)

    qb := r.sb.
        Insert(r.table).
        Columns("session_token", "user_id", "data", "created_at", "expires_at").
        Values(token, userID, sq.Expr(`'{}'::JSONB`), sq.Expr("NOW()"), expires)

    sqlStr, args, err := qb.ToSql()
    if err != nil {
        return "", err
    }
    if _, err := r.db.ExecContext(ctx, sqlStr, args...); err != nil {
        return "", err
    }
    return token, nil
}

// 토큰으로 세션 삭제
func (r *sessionRepo) Delete(ctx context.Context, token string) error {
    qb := r.sb.
        Delete(r.table).
        Where(sq.Eq{"session_token": token})

    sqlStr, args, err := qb.ToSql()
    if err != nil {
        return err
    }
    _, err = r.db.ExecContext(ctx, sqlStr, args...)
    return err
}

// 토큰 유효성 검사 후 세션 반환
func (r *sessionRepo) GetByToken(ctx context.Context, token string) (*model.Session, error) {
    qb := r.sb.
        Select(
            "session_id",
            "session_token",
            "user_id",
            "data",
            "created_at",
            "expires_at",
        ).
        From(r.table).
        Where(sq.Eq{"session_token": token}).
        Where(sq.Gt{"expires_at": sq.Expr("NOW()")})

    sqlStr, args, err := qb.ToSql()
    if err != nil {
        return nil, err
    }

    var s model.Session
    if err := r.db.GetContext(ctx, &s, sqlStr, args...); err != nil {
        return nil, err
    }
    return &s, nil
}
