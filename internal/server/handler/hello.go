package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	router "github.com/yizhinailong/api-demo/internal/server"
)

func init() {
	router.Register(&HelloHandler{})
}

func (h *HelloHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/")
	{
		group.GET("/hello", h.HelloWorld)
		group.GET("/print/:name", h.Print)
	}
}

type HelloHandler struct{}

func (h *HelloHandler) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
		"status":  "success",
	})
}

func (h *HelloHandler) Print(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "print",
		"status":  "success",
		"data":    c.Param("name"),
		"query":   c.Query("name"),
		"header":  c.GetHeader("Authorization"),
		"body":    c.PostForm("name"),
	})
}
