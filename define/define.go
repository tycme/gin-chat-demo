package define

import "os"

var MailPassword = os.Getenv("MailPassword")

type Message struct {
	Message      string `json:"message"`
	RoomIdentity string `json:"room_identity"`
}
