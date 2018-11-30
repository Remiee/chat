package server

import (
	"os"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var storage []Message
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
	join()
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
			leave()
			break
		}

		storage = append(storage, msg)
		saveMessage(msg)

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

func join() {
	var sysMsg = Message{
		Email:    "system",
		Username: "system",
		Message:  "Somebody joined to the channel.",
	}
	broadcast <- sysMsg
}

func leave() {
	var sysMsg = Message{
		Email:    "system",
		Username: "system",
		Message:  "Somebody has left the channel.",
	}
	broadcast <- sysMsg
}
