package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tycme/gin-chat/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/login", service.Login)

	return r
}
