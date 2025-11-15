package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	"github.com/yizhinailong/api-demo/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *bun.DB
	dbOnce sync.Once
)

// GetDB returns the shared database instance, initializing it if necessary
func GetDB() *bun.DB {
	dbOnce.Do(func() {
		cfg := config.GetConfig()

		// Check if database config exists
		if cfg.Database.Driver == "" {
			slog.Error("Database configuration not found")
			return
		}

		// Initialize database connection
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)

		sqldb, err := sql.Open(cfg.Database.Driver, dsn)
		if err != nil {
			slog.Error("Failed to connect to database", "error", err)
			return
		}

		// Test database connection
		if err := sqldb.Ping(); err != nil {
			slog.Error("Failed to ping database", "error", err)
			return
		}

		// Create Bun DB instance
		db = bun.NewDB(sqldb, mysqldialect.New())

		slog.Info("Database connection initialized successfully")
	})

	return db
}
