package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/yizhinailong/api-demo/internal/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	mysqlDB    *bun.DB
	postgresDB *bun.DB
	dbOnce     sync.Once
)

// GetDB returns the shared database instance, initializing it if necessary
func GetMySQLDB() *bun.DB {
	dbOnce.Do(func() {
		cfg := config.GetConfig()

		// Initialize database connection
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.Database.MySQL.User,
			cfg.Database.MySQL.Password,
			cfg.Database.MySQL.Host,
			cfg.Database.MySQL.Port,
			cfg.Database.MySQL.Name,
		)

		sqldb, err := sql.Open("mysql", dsn)
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
		mysqlDB = bun.NewDB(sqldb, mysqldialect.New())

		slog.Info("Database connection initialized successfully")
	})

	return mysqlDB
}
func GetPostgresDB() *bun.DB {
	dbOnce.Do(func() {
		cfg := config.GetConfig()
		// Initialize PostgreSQL database connection
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Database.Postgres.User,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.Name,
		)

		sqldb, err := sql.Open("postgres", dsn)
		if err != nil {
			slog.Error("Failed to connect to PostgreSQL database", "error", err)
			return
		}

		// Test database connection
		if err := sqldb.Ping(); err != nil {
			slog.Error("Failed to ping PostgreSQL database", "error", err)
			return
		}

		// Create Bun DB instance with PostgreSQL dialect
		postgresDB = bun.NewDB(sqldb, pgdialect.New())

		slog.Info("PostgreSQL database connection initialized successfully")
	})

	return postgresDB
}
