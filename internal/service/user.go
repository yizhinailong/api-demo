package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/yizhinailong/api-demo/internal/model"
	"github.com/yizhinailong/api-demo/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	cache    sync.Map // 简单内存缓存（生产用 Redis）
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	// 1. 先查缓存
	if cached, ok := s.cache.Load(id); ok {
		return cached.(*model.User), nil
	}

	// 2. 缓存未命中，查数据库
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存
	s.cache.Store(id, user)

	return user, nil
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	if !strings.Contains(email, "@") {
		return fmt.Errorf("邮箱格式不正确：缺少 @ 符号")
	}

	return nil
}

func (s *UserService) CreateUser(ctx context.Context, input *CreateUserInput) (*model.User, error) {
	// 1. 业务验证：邮箱格式、用户名长度等
	if err := validateEmail(input.Email); err != nil {
		return nil, fmt.Errorf("邮箱验证失败: %w", err)
	}

	if len(input.Username) < 3 {
		return nil, fmt.Errorf("用户名长度至少3个字符")
	}

	// 2. 构造模型
	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
	}

	// 3. 调用 Repository 持久化
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 4. 返回结果（ID 已由 Repository 填充）
	return user, nil
}

// CreateUserInput 创建用户输入（与 model 分离）
type CreateUserInput struct {
	Username string
	Email    string
}
