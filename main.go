package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/enverbisevac/xdefi/internal/httpx"
	"github.com/enverbisevac/xdefi/internal/queue/kafka"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	brokersUrl := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokersUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func(producer sarama.SyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Println(err)
		}
	}(producer)
	consumer, err := kafka.NewConsumer(brokersUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func(consumer sarama.Consumer) {
		err := consumer.Close()
		if err != nil {
			log.Println(err)
		}
	}(consumer)
	queue := kafka.NewQueue(producer, consumer)
	queue.Start(ctx)
	defer queue.Close()

	router := mux.NewRouter()
	server := httpx.NewServer(router, queue)

	go func() {
		log.Fatal(http.ListenAndServe(":8087", server))
	}()

	<-ctx.Done()
	log.Println("Main done")
}
