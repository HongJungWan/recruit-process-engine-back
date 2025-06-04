package model

import "time"

type User struct {
    ID        int       `db:"id" json:"id"`
    Email     string    `db:"email" json:"email"`
    Password  string    `db:"password" json:"-"`
    Name      string    `db:"name" json:"name"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
    return "users"
}
