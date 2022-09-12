package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tycme/gin-chat/define"
	"github.com/tycme/gin-chat/helper"
	"github.com/tycme/gin-chat/models"
)

func Login(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
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

func UserDetail(c *gin.Context) {
	u, _ := c.Get("user_claims")
	uc := u.(*helper.UserClaims)
	userBasic, err := models.GetUserBasicByIdentity(uc.Identity)
	if err != nil {
		log.Printf("[DB ERROR]: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": userBasic,
	})
}

func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	cnt, err := models.GetUserBasicCountByEmail(email)
	if err != nil {
		log.Printf("[DB ERROR: %v\n]", err)
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已被注册",
		})
		return
	}
	code := helper.GetCode()
	err = helper.SendCode(email, code)
	if err != nil {
		log.Printf("[ERROR: %v]\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "email不能为空",
		})
		return
	}
	if err := models.RDB.Set(context.Background(), define.RegisterPrefix+email, code, time.Second*time.Duration(define.ExpireTime)).Err(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
	})
}

func Regiseter(c *gin.Context) {
	code := c.PostForm("code")
	email := c.PostForm("email")
	account := c.PostForm("acconut")
	password := c.PostForm("password")
	if code == "" || email == "" || account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	// 判断账号是否唯一
	cnt, err := models.GetUserBasicCountByEmail(email)
	if err != nil {
		log.Printf("[ERROR] : %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号已被注册",
		})
		return
	}
	// 验证码是否正确
	r, err := models.RDB.Get(context.Background(), define.RegisterPrefix+email).Result()
	if err != nil || r != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确",
		})
		return
	}
	ub := &models.UserBasic{
		Identity:  helper.GetUuid(),
		Account:   account,
		Password:  helper.GetMd5(password),
		Nickname:  "",
		Sex:       0,
		Email:     email,
		CreatedAt: time.Now().Unix(),
		UpdateAt:  time.Now().Unix(),
	}
	err = models.InsertOneUserBasic(ub)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
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

type UserQueryResult struct {
	Nickname string `json:"nickname"`
	Sex      int    `bson:"sex"`
	Email    string `bson:"email"`
	Avatar   string `bson:"avatar"`
	IsFriend bool   `json:"is_friend"`
}

func UserQuery(c *gin.Context) {
	account := c.Query("account")
	userBasic, err := models.GetUserBasicByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
	}
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
	}
	uc, _ := c.MustGet("user_claims").(*helper.UserClaims)
	isFriend := models.JudgeUserIsFriend(userBasic.Identity, uc.Identity)
	uq := UserQueryResult{
		Nickname: userBasic.Nickname,
		Sex:      userBasic.Sex,
		Email:    userBasic.Email,
		Avatar:   userBasic.Avatar,
		IsFriend: isFriend,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": uq,
	})
}

func UserAdd(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	ub, err := models.GetUserBasicByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "没有该用户",
		})
		return
	}
	uc, _ := c.MustGet("user_claims").(*helper.UserClaims)
	if models.JudgeUserIsFriend(uc.Identity, ub.Identity) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "已经是好友,不可重复添加",
		})
		return
	}
	// 保存房间记录
	rb := &models.RoomBasic{
		Identity:     helper.GetUuid(),
		UserIdentity: uc.Identity,
		CreatedAt:    time.Now().Unix(),
		UpdateAt:     time.Now().Unix(),
	}
	err = models.InsertOneRoomBasic(rb)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	// 保存用户与房间关联关系
	ur := models.UserRoom{
		UserIdentity: ub.Identity,
		RoomIdentity: uc.Identity,
		RoomType:     1,
		CreatedAt:    time.Now().Unix(),
		UpdateAt:     time.Now().Unix(),
	}
	if err := models.InsertOneUserRoom(&ur); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "好友添加成功",
	})
}
