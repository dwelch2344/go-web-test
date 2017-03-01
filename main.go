package main

import (
	"github.com/derekdowling/go-json-spec-handler/jsh-api"
	"github.com/Shopify/sarama"

	"net/http"
	"log"
	"flag"

	"github.com/aiur/web-test/ws"
	"github.com/aiur/web-test/kafka"
)

func main() {

	read := flag.Bool("c", false, "Create a consumer")
	write := flag.Bool("p", false, "Create a producer")
	flag.Parse()


	//addresses of available kafka brokers
	brokers := []string{"localhost:9092"}
	//setup relevant config info
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	if *read {
		log.Println("Creating consumer")
		kafka.DoKafkaRead(brokers, config)
	}else if *write {
		log.Println("Creating producer")
		kafka.DoKafkaWrite(brokers, config)
	}else {
		log.Println("Run webserver")

		master, err := sarama.NewConsumer(brokers, config)
		if err != nil {
			panic(err)
		}

		topic := "dave1"
		consumer, err := master.ConsumePartition(topic, int32(0), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		ws.Run(consumer)
	}
}




func doConsul(){
	userStorage := &UserStorage{}
	resource := jshapi.NewCRUDResource("user", userStorage)
	//resource.UseC(yourUserMiddleware)

	// setup a logger, your shiny new API, and give it a resource
	//logger := log.New(os.Stderr, "<yourapi>: ", log.LstdFlags)
	//api := jshapi.Default("/test", true, logger)
	api := jshapi.Default("/test", true, nil)
	api.Add(resource)

	// launch your api
	http.ListenAndServe("localhost:8000", api)
}


