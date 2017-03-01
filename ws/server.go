package ws


import (
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
	"github.com/Shopify/sarama"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func Run(consumer sarama.PartitionConsumer){
	r := mux.NewRouter()
	r.HandleFunc("/ws", curry(consumer))
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}


type strategy func(w http.ResponseWriter, r *http.Request) ()



func curry(consumer sarama.PartitionConsumer) strategy {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming websocket connection")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Upgraded connection")

		signals := make(chan bool, 1)
		conn.SetCloseHandler(func(code int, text string) error {
			log.Println("Closing", code, text)
			signals <- true
			return nil
		})

		// Count how many message processed
		msgCount := 0

		// Get signnal for finish
		doneCh := make(chan bool)
		go doFmt(consumer, &msgCount, conn, doneCh, signals)

		<-doneCh
		fmt.Println("Processed", msgCount, "messages")

	}
}


func doFmt(consumer sarama.PartitionConsumer, msgCount *int,
		   conn *websocket.Conn, doneCh chan<- bool,
	       signals <-chan bool) {


	for {
		select {
		case err := <-consumer.Errors():
			fmt.Println(err)
		case msg := <-consumer.Messages():
			*msgCount++
			fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			if err := conn.WriteMessage(websocket.TextMessage, msg.Value); err != nil {
				log.Println("Exiting from err on write:", err)
				return
			}
		case <-signals:
			fmt.Println("Connection closed")
			doneCh <- true
		}
	}
}