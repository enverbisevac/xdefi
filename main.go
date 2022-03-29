package main

import (
	"context"
	"github.com/enverbisevac/xdefi/internal/httpx"
	"github.com/enverbisevac/xdefi/internal/queue/kafka"
	"github.com/enverbisevac/xdefi/internal/service"
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	brokersUrl := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokersUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	consumer, err := kafka.NewConsumer(brokersUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	queue := kafka.NewQueue(producer, consumer)

	router := mux.NewRouter()
	server := httpx.NewServer(ctx, router, queue)

	signer := service.NewSigner(queue)
	signer.Start(ctx)

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-ctx.Done()
	log.Println("Main done")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}
