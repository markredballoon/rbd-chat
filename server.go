package main

import (
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {

		log.Println("on connection")
		so.Join("chat")

		so.On("chat message", func(msg string) {
			// Build return message
			returnMsg := map[string]string{"title": "Message send on ", "content": msg}

			// Add timestamp
			now := time.Now()
			nowFormatted := now.Format("3:01pm")

			returnMsg["title"] += nowFormatted
			returnEncoded, _ := json.Marshal(returnMsg)

			log.Printf("returned: %s \n", returnEncoded)

			so.Emit("chat message", string(returnEncoded))
			so.BroadcastTo("chat", "chat message", string(returnEncoded))
		})

		so.On("chat message with ack", func(msg string) string {
			log.Printf("received: %s \n", msg)
			return "Message Reveived"
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
