package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tycme/gin-chat/middlewares"
	"github.com/tycme/gin-chat/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/login", service.Login)
	r.POST("/send/code", service.SendCode)
	r.POST("/register", service.Regiseter)

	auth := r.Group("/u", middlewares.AuthCheck())
	auth.GET("/user/detail", service.UserDetail)
	auth.GET("/user/query", service.UserQuery)

	auth.GET("/websocket/message", service.WebsocketMessage)
	auth.GET("/chat/list", service.ChatList)

	auth.POST("/user/add", service.UserAdd)
	auth.DELETE("/user/delete", service.Delete)
	return r
}
