package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type BidEvent struct {
	Price int
}

var (
	currentPrice int

	eventChan = make(chan BidEvent, 100)

	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func broadcast(price int) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		_ = conn.WriteJSON(map[string]int{
			"price": price,
		})
	}
}

func run() {
	for event := range eventChan {

		if event.Price > currentPrice {
			currentPrice = event.Price
			fmt.Println("New Price:", currentPrice)

			broadcast(currentPrice)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	for {
		var msg map[string]string

		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		price, _ := strconv.Atoi(msg["price"])

		eventChan <- BidEvent{
			Price: price,
		}
	}
}


func main(){
	go run()
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}