package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("err: %v", err)
			delete(clients, ws)
			break
		}

		ObsceneFilter(&msg)
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

func ObsceneFilter(msg *Message) {
	obscenes := [6]string{"buzi", "fasz", "szar", "anyád", "kurva", "geci"}
	var result []string
	for _, word := range strings.Split(msg.Message, " ") {
		for _, obsceneWord := range obscenes {
			if strings.Contains(word, obsceneWord) {
				word = "****"
			}
		}
		result = append(result, word)
	}
	msg.Message = strings.Join(result, " ")
}
