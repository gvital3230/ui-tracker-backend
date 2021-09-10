package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"ui-tracker-backend/internal/ws"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws-public", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r, false)
	})
	http.HandleFunc("/ws-dashboard", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r, true)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
