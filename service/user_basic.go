package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tycme/gin-chat/helper"
	"github.com/tycme/gin-chat/models"
)

func Login(c *gin.Context) {
	account := c.PostForm("account")
	fmt.Println("account: ", account)
	password := c.PostForm("password")
	fmt.Println("password: ", password)
	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	ub, err := models.GetUserBasicByAccountPassword(account, helper.GetMd5(password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}
	token, err := helper.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误：" + err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登陆成功",
		"data": gin.H{
			"token": token,
		},
	})
}
