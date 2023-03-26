package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	runWS()
}

func runWS() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("starting ws server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP request to WebSocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send a message to the client
	message := []byte("Hello, client!")
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for incoming messages from the client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Received message: %s\n", p)

		// Respond to the message
		response := []byte("You said: " + string(p))
		err = conn.WriteMessage(messageType, response)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
