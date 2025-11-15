package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yizhinailong/api-demo/internal/repository"
	"github.com/yizhinailong/api-demo/internal/service"

	_ "github.com/go-sql-driver/mysql"

	router "github.com/yizhinailong/api-demo/internal/server"
)

func init() {
	// Initialize repository and service
	userRepo := repository.NewUserMySQLRepository(repository.GetDB())
	userService := service.NewUserService(userRepo)

	// Register handler with initialized service
	router.Register(&UserHandler{userService: userService})
}

type UserHandler struct {
	userService *service.UserService
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/users")
	{
		group.POST("/create", h.CreateUser)
		group.GET("/:id", h.GetUser)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input service.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.userService.CreateUser(c, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	user, err := h.userService.GetUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
