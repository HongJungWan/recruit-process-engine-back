package repository

import (
	model "github.com/HongJungWan/recruit-process-engine-back/internal/user/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
}

type userRepo struct {
    db     *sqlx.DB
    sb     sq.StatementBuilderType
    table  string
}

// 의존성 주입(DI)
func NewUserRepository(db *sqlx.DB) UserRepository {
    tableName := (&model.User{}).TableName()
    
    return &userRepo{
        db:    db,
        sb:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
        table: tableName,
    }
}

