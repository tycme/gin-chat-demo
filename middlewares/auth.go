package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tycme/gin-chat/helper"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		userClaims, err := helper.AnalyseToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户认证失败",
			})
			return
		}
		c.Set("user_claims", userClaims)
		c.Next()
	}
}
