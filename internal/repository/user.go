package repository

import (
	"context"

	"github.com/yizhinailong/api-demo/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error // 注意使用 model.User
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*model.User, error)
}
