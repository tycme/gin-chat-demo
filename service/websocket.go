package service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tycme/gin-chat/define"
	"github.com/tycme/gin-chat/helper"
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
		for _, cc := range wc {
			err := cc.WriteMessage(websocket.TextMessage, []byte(msg.Message))
			if err != nil {
				log.Printf("cc.WriteMessage: %v\n", err)
				return
			}
		}
	}
}
