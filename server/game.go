package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var blue bool = true

var c *websocket.Conn

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q", html.EscapeString((r.URL.Path)))
}

func toggle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	blue = !blue
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, `{"blue": %t}`, blue)
}

var upgrader = websocket.Upgrader{}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
}

func main() {

	http.HandleFunc("/", home)

	http.HandleFunc("/toggle", toggle)

	http.HandleFunc("/get", get)

	http.HandleFunc("/ws", serveWs)

	log.Fatal(http.ListenAndServe(":8081", nil))

}
