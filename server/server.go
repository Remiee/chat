package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var storage []Event
var clients = make(map[*websocket.Conn]User)
var broadcast = make(chan Event)
var upgrader = websocket.Upgrader{}

type Event struct {
	ID      string    `json:"id"`
	Type    EventType `json:"type"`
	User    User      `json:"user"`
	Message string    `json:"message"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type EventType int8

const (
	Msg   EventType = 0
	Join  EventType = 1
	Leave EventType = 2
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()
	clients[ws] = User{}

	for _, msg := range storage {
		ws.WriteJSON(msg)
	}

	for {
		var msg Event
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("err: %v", err)
			leave(clients[ws])
			delete(clients, ws)
			break
		}

		if msg.Type == Join {
			join(msg.User)
			clients[ws] = msg.User
			continue
		}

		// msg.EventType = 0

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

func saveMessage(msg Event) {
	f, err := os.OpenFile("./message.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(msg.User.Username + "/" + msg.User.Email + " - " + msg.Message + "\n"))
	if err != nil {
		panic(err)
	}
	f.Close()
}

func ObsceneFilter(msg *Event) {
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

func join(user User) {
	var sysMsg = Event{
		User: User{
			Username: "system",
			Email:    "system",
		},
		Message: user.Username + " joined to the channel.",
	}
	broadcast <- sysMsg
}

func leave(user User) {
	var sysMsg = Event{
		User: User{
			Username: "system",
			Email:    "system",
		},
		Message: user.Username + " has left the channel.",
	}
	broadcast <- sysMsg
}
