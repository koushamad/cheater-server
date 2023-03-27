package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const Client = "server"

// Message struct
type Message struct {
	ApiKey  string `json:"apiKey"`
	Client  string `json:"client"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240,
	WriteBufferSize: 10240,
}

type client struct {
	conn   *websocket.Conn
	apiKey string
	client string
}

var clients = make(map[*client]bool)
var broadcast = make(chan *Message, 8)

func main() {
	// Upgrade HTTP connections to WebSocket connections
	http.HandleFunc("/ws", handleWebSocket)

	// Start broadcasting messages
	go handleMessages()

	// Listen on port 8080
	log.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}

	// Read the ID and client from the first message
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error: ", err)
		return
	}

	var message Message
	err = json.Unmarshal(msg, &message)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return
	}

	// Create a new client
	c := &client{conn: conn, apiKey: message.ApiKey, client: message.Client}
	clients[c] = true

	// Listen for messages from the client
	go handleClientMessages(c)
}

func handleClientMessages(c *client) {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Read error: ", err)
			delete(clients, c)
			return
		}

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Unmarshal error: ", err)
			continue
		}

		// Broadcast the message to all clients with the same ID
		broadcast <- &message
	}
}

func handleMessages() {
	for {
		// Get the next message from the broadcast channel
		message := *<-broadcast

		// Send the message to all clients with the same ID
		if message.Client != Client {
			sendToClient(message)
		}
	}
}

func sendToClient(message Message) {
	for c := range clients {
		if c.apiKey == message.ApiKey && c.client == message.Client {
			err := c.conn.WriteJSON(message)
			if err != nil {
				log.Println("Write error: ", err)
				delete(clients, c)

			}
		}
	}
}
