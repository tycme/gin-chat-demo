package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tycme/gin-chat/middlewares"
	"github.com/tycme/gin-chat/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/login", service.Login)

	auth := r.Group("/u", middlewares.AuthCheck())
	auth.GET("/user/detail", service.UserDetail)

	return r
}
