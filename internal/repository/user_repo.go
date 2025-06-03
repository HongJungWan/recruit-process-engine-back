package repository

import (
    "context"
    "github.com/jmoiron/sqlx"
    sq "github.com/Masterminds/squirrel" // Squirrel import
    "github.com/HongJungWan/recruit-process-engine-back/internal/models"
)

type UserRepository interface {
    GetByID(ctx context.Context, id int) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Create(ctx context.Context, user *models.User) (int, error)
}

type userRepo struct {
    db     *sqlx.DB
    sb     sq.StatementBuilderType
    table  string
}

// 의존성 주입(DI) 패턴
func NewUserRepository(db *sqlx.DB) UserRepository {
    tableName := (&models.User{}).TableName()
    return &userRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: tableName,
    }
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
    queryBuilder := r.sb.
        Select("id", "email", "password", "name", "created_at", "updated_at").
        From(r.table).
        Where(sq.Eq{"id": id})

    sqlStr, args, err := queryBuilder.ToSql()
    if err != nil {
        return nil, err
    }

    var u models.User
    if err := r.db.GetContext(ctx, &u, sqlStr, args...); err != nil {
        return nil, err
    }

    return &u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    queryBuilder := r.sb.
        Select("id", "email", "password", "name", "created_at", "updated_at").
        From(r.table).
        Where(sq.Eq{"email": email})

    sqlStr, args, err := queryBuilder.ToSql()
    if err != nil {
        return nil, err
    }

    var u models.User
    if err := r.db.GetContext(ctx, &u, sqlStr, args...); err != nil {
        return nil, err
    }

    return &u, nil
}

func (r *userRepo) Create(ctx context.Context, user *models.User) (int, error) {
    queryBuilder := r.sb.
        Insert(r.table).
        Columns("email", "password", "name", "created_at", "updated_at").
        Values(user.Email, user.Password, user.Name, sq.Expr("NOW()"), sq.Expr("NOW()")).
        Suffix("RETURNING id")

    sqlStr, args, err := queryBuilder.ToSql()
    if err != nil {
        return 0, err
    }

    var newID int
    if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&newID); err != nil {
        return 0, err
    }

    return newID, nil
}
