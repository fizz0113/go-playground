package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	log.SetPrefix("[client] ")
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("received from server: %s\n", msg)
		}
	}()

	for {
		time.Sleep(time.Second)
		err := c.WriteMessage(websocket.TextMessage, []byte("Hello Server"))
		if err != nil {
			log.Println("write error:", err)
			return
		}
	}
}
