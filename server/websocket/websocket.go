package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redisStream/redis"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{} // use default options
var clientMsgChan = make(chan bool)

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ERROR Upgrade:", err)
	} else {
		fmt.Println("client connected")
	}

	return ws, err
}

func Ping(ws *websocket.Conn) {
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		} else {
			redis_persistence.GetRedis().PopHeadData()
			clientMsgChan <- true
			log.Printf("recv: %s", message)
		}
	}
}

func ClientWriter(ws *websocket.Conn) {
	defer ws.Close()
	for {
		msg := redis_persistence.GetRedis().ReadHeadData()
		if msg != "" {
			log.Println("Sending ", msg)
			err := ws.WriteMessage(1, []byte(msg))
			if err != nil {
				log.Println("write:", err)
				return
			}
			<-clientMsgChan
		}
		time.Sleep(time.Second)
	}
}
