package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
	"os"
	"os/signal"
)

func DoKafkaRead(brokers []string, config *sarama.Config) {

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	topic := "dave1"
	//consumer, err := master.ConsumePartition(topic, int32(1), sarama.OffsetOldest)
	consumer, err := master.ConsumePartition(topic, int32(0), sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
}

func DoKafkaWrite(brokers []string, config *sarama.Config){

	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		log.Fatalf("Failed 1 – %v", err)
	}



	topic := "dave1" //e.g create-user-topic
	partition := int32(1) //Partition to produce to

	ticker := time.NewTicker(time.Second * 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	msgCount := 0
	now := time.Now()

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:

				msg := fmt.Sprintf("%v-%v", now.Format("15:04:05"), msgCount)
				message := &sarama.ProducerMessage{
					Topic: topic,
					Partition: partition,
					Value: sarama.StringEncoder(msg),
				}

				partition, offset, err := producer.SendMessage(message)

				if err != nil {
					log.Fatalf("Failed sending message – %v", err)
				}

				log.Printf("Printed to %v - %v", partition, offset)
				msgCount++
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Sent", msgCount, "messages")
}