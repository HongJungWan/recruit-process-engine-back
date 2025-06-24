package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/user/model"
)

type UserRepository interface {
    GetByEmail(ctx context.Context, loginID string) (*model.User, error)
    Authenticate(ctx context.Context, loginID, password string) (int, error)
}

type userRepo struct {
    db    *sqlx.DB
    sb    sq.StatementBuilderType
    table string
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: (&model.User{}).TableName(),
    }
}

// 로그인 ID(login_id)로 유저 조회
func (r *userRepo) GetByEmail(ctx context.Context, loginID string) (*model.User, error) {
    qb := r.sb.
        Select(
            "user_id",
            "login_id",
            "login_pw",
            "name",
            "email",
            "role",
            "created_at",
            "created_by",
            "updated_at",
            "updated_by",
        ).
        From(r.table).
        Where(sq.Eq{"login_id": loginID})

    sqlStr, args, err := qb.ToSql()
    if err != nil {
        return nil, err
    }

    var u model.User
    if err := r.db.GetContext(ctx, &u, sqlStr, args...); err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *userRepo) Authenticate(ctx context.Context, loginID, password string) (int, error) {
    // 사용자 정보 조회
    u, err := r.GetByEmail(ctx, loginID)
    if err != nil {
        return 0, err
    }

    // TODO: 비밀번호 비교

    return u.UserID, nil
}