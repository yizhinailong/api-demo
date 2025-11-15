package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64     `bun:",pk,autoincrement" json:"id"`
	Username  string    `bun:"username,notnull" json:"username"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
