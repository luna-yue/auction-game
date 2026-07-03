package main

import (
	"fmt"
	"net/http"
)

var currentPrice int
var eventChan = make(chan int, 10)

func run(){
	for price:=range eventChan{
		if price > currentPrice{
			currentPrice = price
			fmt.Println("Current Price Updated:", currentPrice)
		}
	}
}

func bid(w http.ResponseWriter, r *http.Request){
	price := r.URL.Query().Get("price")
	if price == ""{
		http.Error(w, "Price is required", http.StatusBadRequest)
		return
	}

	var bidPrice int
	_, err := fmt.Sscanf(price, "%d", &bidPrice)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	eventChan <- bidPrice
	fmt.Fprintf(w, "Bid received: %d", bidPrice)
}


func main(){
	go run()
	http.HandleFunc("/bid", bid)
	http.ListenAndServe(":8080", nil)
}