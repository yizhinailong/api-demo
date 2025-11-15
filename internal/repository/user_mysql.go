package repository

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/yizhinailong/api-demo/internal/model"
)

type userMySQLRepo struct {
	db *bun.DB
}

func NewUserMySQLRepository(db *bun.DB) UserRepository {
	return &userMySQLRepo{db: db}
}

func (r *userMySQLRepo) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

func (r *userMySQLRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

func (r *userMySQLRepo) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

func (r *userMySQLRepo) Delete(ctx context.Context, id int64) error {
	user := &model.User{ID: id}
	_, err := r.db.NewDelete().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}
	return nil
}

func (r *userMySQLRepo) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	return users, nil
}
