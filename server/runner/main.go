package main

import (
	"flag"
	"fmt"
	"github.com/redisStream/redis"
	"github.com/redisStream/server/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hmmmmmmm")
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	go websocket.ClientWriter(ws)
	websocket.Ping(ws)
}

func setupRoutes() {
	http.HandleFunc("/ping", ping)
}

func insertData() {
	go func() {
		count := 0
		for {
			msg := strconv.Itoa(count)
			redis_persistence.GetRedis().InsertTailData(msg)
			count ++
			time.Sleep(5 * time.Second)
		}
	}()
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	setupRoutes()
	insertData()
	addr := flag.String("addr", "localhost:8080", "http service address")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
