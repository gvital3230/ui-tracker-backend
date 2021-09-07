package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"html"
	"log"
	"net/http"
	"os"
	"strings"
)

func echo(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
			return strings.Contains(allowedOrigins, origin)
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/echo", echo)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
