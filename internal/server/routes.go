package router

import "github.com/gin-gonic/gin"

type RouteRegistrar interface {
	RegisterRoutes(r *gin.Engine)
}

var registrars []RouteRegistrar

func Register(registrar RouteRegistrar) {
	registrars = append(registrars, registrar)
}

func SetupRoutes(r *gin.Engine) {
	for _, reg := range registrars {
		reg.RegisterRoutes(r)
	}
}
