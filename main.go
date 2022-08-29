package main

import (
	"github.com/tycme/gin-chat/router"
)

func main() {
	e := router.Router()
	e.Run(":8081")
}
