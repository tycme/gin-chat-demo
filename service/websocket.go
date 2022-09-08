package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tycme/gin-chat/define"
	"github.com/tycme/gin-chat/helper"
	"github.com/tycme/gin-chat/models"
)

var upgrader = websocket.Upgrader{}

var wc = make(map[string]*websocket.Conn)

func WebsocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常， upgrade websocket失败, err: " + err.Error(),
		})
		return
	}
	defer conn.Close()
	uc := c.MustGet("user_claims").(*helper.UserClaims)
	wc[uc.Identity] = conn
	for {
		msg := &define.Message{}
		err := conn.ReadJSON(msg)
		if err != nil {
			log.Printf("read error: %v\n", err)
			return
		}
		// 判断用户是否属于消息体的房间
		_, err = models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, msg.RoomIdentity)
		if err != nil {
			log.Printf("没找到user room数据")
			return
		}
		// 保存消息
		msgBasic := &models.MsgBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: msg.RoomIdentity,
			Data:         msg.Message,
			CreatedAt:    time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}
		err = models.InsertOneMsg(msgBasic)
		if err != nil {
			log.Printf("[DB ERROR: %+v]\n", err)
		}
		// 获取特定房间在线用户
		userRooms, err := models.GetUserRoomByRoomIdentity(msg.RoomIdentity)
		if err != nil {
			log.Printf("[DB ERROR]: %+v", err)
		}
		for _, room := range userRooms {
			if cc, ok := wc[room.UserIdentity]; ok {
				err := cc.WriteMessage(websocket.TextMessage, []byte(msg.Message))
				if err != nil {
					log.Printf("cc.WriteMessage: %v\n", err)
					return
				}
			}
		}
	}
}
