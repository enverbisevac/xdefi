package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

type Queue struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewQueue(producer sarama.SyncProducer, consumer sarama.Consumer) Queue {
	q := Queue{
		producer: producer,
		consumer: consumer,
	}
	return q
}

func (q Queue) PushToQueue(topic string, message []byte) error {
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

func (q Queue) Register(ctx context.Context, topic string) chan string {
	resultChan := make(chan string, 100)
	consumer, err := q.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	log.Printf("Consumer started with topic: %s", topic)
	go func() {
		defer close(resultChan)
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				log.Printf("Received message: | Topic(%s) | Message(%s) \n", msg.Topic, string(msg.Value))
				resultChan <- string(msg.Value)
			case <-ctx.Done():
				fmt.Println("Context done")
				return
			}
		}
	}()
	return resultChan
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
