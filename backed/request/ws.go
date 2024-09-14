package request

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"main/binary"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WebSocketHandle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		binary.DebugLog.Println(err)
		return
	}
	userid := c.Query("name")
	ip := c.ClientIP()
	id := ip + userid
	idMd5 := fmt.Sprintf("%x", md5.Sum([]byte(id)))
	client := &Client{
		ID:         idMd5,
		Socket:     conn,
		Send:       make(chan []byte, 10),
		UserId:     userid,
		Start:      time.Now(),
		ExpireTime: time.Minute * 10,
	}
	Manager.Register <- client
	go client.Write()
	go client.Check()
	go client.Read()
}
func Transform(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		userInput := string(msg)

		err = ws.WriteMessage(websocket.TextMessage, []byte(userInput+" txt"))
		if err != nil {
			log.Println("Write back", err)
			break
		}
	}
}
