package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yizhinailong/api-demo/internal/service"

	_ "github.com/go-sql-driver/mysql"

	router "github.com/yizhinailong/api-demo/internal/server"
)

func init() {
	// Initialize repository and service
	userService := service.NewUserService()

	// Register handler with initialized service
	router.Register(&UserHandler{userService: userService})
}

type UserHandler struct {
	userService service.UserServiceInterface
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/users")
	{
		group.POST("/create", h.CreateUser)
		group.GET("/get", h.GetUser)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input service.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, err := h.userService.CreateUser(c, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "user created", "userID": user.ID})
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var input service.GetUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, err := h.userService.GetUser(c, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
	}
}
