package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/yizhinailong/api-demo/internal/config"

	_ "github.com/yizhinailong/api-demo/internal/server/handler"

	router "github.com/yizhinailong/api-demo/internal/server"
)

func main() {
	cfg := config.GetConfig()

	engine := gin.Default()
	router.SetupRoutes(engine)

	slog.Info("Web server listening address http://localhost:" + cfg.Server.Port)

	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		slog.Error("Error starting server error " + err.Error())
	}
}
