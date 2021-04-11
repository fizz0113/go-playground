package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 升級 HTTP 連線為 WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允許跨網域（開發用）
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	// 將連線升級成 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("client connected")

	for {
		// 接收訊息
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("received: %s\n", message)

		// 回傳訊息給 client
		err = conn.WriteMessage(mt, []byte("server got: "+string(message)))
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func main() {
	log.SetPrefix("[server] ")
	http.HandleFunc("/ws", handleWS)
	log.Println("webSocket server started on :8080")
	http.ListenAndServe(":8080", nil)
}
