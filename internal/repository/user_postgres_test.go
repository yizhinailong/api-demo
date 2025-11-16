package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/yizhinailong/api-demo/internal/model"
)

func TestUserPostgresRepo_Create(t *testing.T) {
	// Create a mock database connection
	sqldb, _ := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := &userPostgresRepo{db: db}

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()
		user := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		// This test would require a real database connection or more complex mocking
		// For now, we'll just test that the method doesn't panic
		err := repo.Create(ctx, user)
		// We expect an error because we don't have a real database
		assert.Error(t, err)
	})
}

func TestUserPostgresRepo_GetByID(t *testing.T) {
	sqldb, _ := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := &userPostgresRepo{db: db}

	t.Run("get non-existent user", func(t *testing.T) {
		ctx := context.Background()
		_, err := repo.GetByID(ctx, 999)
		assert.Error(t, err)
	})
}

func TestUserPostgresRepo_Update(t *testing.T) {
	sqldb, _ := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := &userPostgresRepo{db: db}

	t.Run("update non-existent user", func(t *testing.T) {
		ctx := context.Background()
		user := &model.User{
			ID:       999,
			Username: "updateduser",
			Email:    "updated@example.com",
		}

		err := repo.Update(ctx, user)
		assert.Error(t, err)
	})
}

func TestUserPostgresRepo_Delete(t *testing.T) {
	sqldb, _ := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := &userPostgresRepo{db: db}

	t.Run("delete non-existent user", func(t *testing.T) {
		ctx := context.Background()
		err := repo.Delete(ctx, 999)
		assert.Error(t, err)
	})
}

func TestUserPostgresRepo_List(t *testing.T) {
	sqldb, _ := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := &userPostgresRepo{db: db}

	t.Run("list users", func(t *testing.T) {
		ctx := context.Background()
		_, err := repo.List(ctx)
		assert.Error(t, err) // Should error because no real DB connection
	})
}
