package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yizhinailong/api-demo/internal/model"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.User), args.Error(1)
}

func TestNewUserService(t *testing.T) {
	service := NewUserService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.mysqlRepo)
	assert.NotNil(t, service.postgresRepo)
}

func TestUserService_getUserRepo(t *testing.T) {
	service := NewUserService()

	t.Run("get mysql repo", func(t *testing.T) {
		repo := service.getUserRepo("mysql")
		assert.NotNil(t, repo)
	})

	t.Run("get postgres repo", func(t *testing.T) {
		repo := service.getUserRepo("postgres")
		assert.NotNil(t, repo)
	})

	t.Run("get default repo", func(t *testing.T) {
		repo := service.getUserRepo("")
		assert.NotNil(t, repo)
	})
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := &UserService{
		mysqlRepo:    mockRepo,
		postgresRepo: mockRepo,
	}

	t.Run("get user from cache", func(t *testing.T) {
		ctx := context.Background()
		expectedUser := &model.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		// Pre-populate cache
		service.cache.Store(int64(1), expectedUser)

		input := &GetUserInput{
			ID:       1,
			Database: "mysql",
		}

		user, err := service.GetUser(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("get user from database", func(t *testing.T) {
		ctx := context.Background()
		expectedUser := &model.User{
			ID:       2,
			Username: "testuser2",
			Email:    "test2@example.com",
		}

		mockRepo.On("GetByID", ctx, int64(2)).Return(expectedUser, nil)

		input := &GetUserInput{
			ID:       2,
			Database: "mysql",
		}

		user, err := service.GetUser(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, assert.AnError)

		input := &GetUserInput{
			ID:       999,
			Database: "mysql",
		}

		_, err := service.GetUser(ctx, input)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestValidateEmail(t *testing.T) {
	t.Run("valid email", func(t *testing.T) {
		err := validateEmail("test@example.com")
		assert.NoError(t, err)
	})

	t.Run("empty email", func(t *testing.T) {
		err := validateEmail("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "邮箱不能为空")
	})

	t.Run("email without @", func(t *testing.T) {
		err := validateEmail("invalid-email")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "邮箱格式不正确")
	})
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := &UserService{
		mysqlRepo:    mockRepo,
		postgresRepo: mockRepo,
	}

	t.Run("create user successfully", func(t *testing.T) {
		ctx := context.Background()
		input := &CreateUserInput{
			Username: "testuser",
			Email:    "test@example.com",
			Database: "mysql",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)

		user, err := service.CreateUser(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, input.Username, user.Username)
		assert.Equal(t, input.Email, user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create user with invalid email", func(t *testing.T) {
		ctx := context.Background()
		input := &CreateUserInput{
			Username: "testuser",
			Email:    "invalid-email",
			Database: "mysql",
		}

		_, err := service.CreateUser(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "邮箱验证失败")
	})

	t.Run("create user with short username", func(t *testing.T) {
		ctx := context.Background()
		input := &CreateUserInput{
			Username: "ab",
			Email:    "test@example.com",
			Database: "mysql",
		}

		_, err := service.CreateUser(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "用户名长度至少3个字符")
	})

	t.Run("create user with database error", func(t *testing.T) {
		ctx := context.Background()
		input := &CreateUserInput{
			Username: "testuser",
			Email:    "test@example.com",
			Database: "mysql",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(assert.AnError)

		_, err := service.CreateUser(ctx, input)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
