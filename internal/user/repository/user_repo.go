package repository

import (
	// 표준 라이브러리
	"context"
	"log"

	// 서드파티(외부) 라이브러리
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/user/model"
)

type UserRepository interface {
    Authenticate(ctx context.Context, loginID, password string) (int, error)
    GetByUserId(ctx context.Context, loginID string) (*model.User, error)
    GetByPassword(ctx context.Context, password string) (*model.User, error)
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

func (r *userRepo) Authenticate(ctx context.Context, loginID, password string) (int, error) {
    // 유저 ID 조회
    ue, err := r.GetByUserId(ctx, loginID)
    if err != nil {
        return 0, err
    }
    log.Printf("[Authenticate] 유저 ID: %d\n", ue.UserID)
    log.Printf("[Authenticate] 로그인 ID: %s\n", ue.LoginID)

    // 유저 비밀번호 조회
    up, err := r.GetByPassword(ctx, password)
    if err != nil {
        return 0, err
    }
    log.Printf("[Authenticate] 비밀번호: %s\n", up.Password)

    return ue.UserID, nil
}

// 유저 ID 조회
func (r *userRepo) GetByUserId(ctx context.Context, loginID string) (*model.User, error) {
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

// 유저 패스워드 조회
func (r *userRepo) GetByPassword(ctx context.Context, password string) (*model.User, error) {
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
        Where(sq.Eq{"login_pw": password})

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
