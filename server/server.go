package server

import (
	"os"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var storage []Message
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for _, msg := range storage {
		ws.WriteJSON(msg)
	}

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("err: %v", err)
			delete(clients, ws)
			break
		}

		storage = append(storage, msg)
		saveMessage(msg)
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func saveMessage(msg Message) {
	f, err := os.OpenFile("./message.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(msg.Username + "/" + msg.Email + " - " + msg.Message + "\n"))
	if err != nil {
        panic(err)
    }
    f.Close()
}
