package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

const (
	SignEvent = "SignEvent"
)

type Queue struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	done     chan struct{}
}

func NewQueue(producer sarama.SyncProducer, consumer sarama.Consumer) Queue {
	q := Queue{
		producer: producer,
		consumer: consumer,
		done:     make(chan struct{}),
	}
	return q
}

func (q Queue) PushToQueue(message []byte) error {
	topic := SignEvent
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := q.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func (q Queue) Start(ctx context.Context) {
	consumer, err := q.consumer.ConsumePartition(SignEvent, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	log.Println("Consumer started ")
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				log.Printf("Received message: | Topic(%s) | Message(%s) \n", msg.Topic, string(msg.Value))
			case <-q.done:
				fmt.Println("Consumer done")
				return
			case <-ctx.Done():
				fmt.Println("Context done")
				return
			}
		}
	}()
}

func (q Queue) Close() {
	q.done <- struct{}{}
}

func NewProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	// NewConsumer creates a new consumer using the given broker addresses and configuration
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
