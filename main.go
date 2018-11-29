package main

import (
	"log"
	"net/http"
	"github.com/Remiee/chat/server"
)

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", server.HandleConnections)
	go server.HandleMessages()

	log.Println("http server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
