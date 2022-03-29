package queue

import "context"

const (
	SignEvent   = "SignEvent"
	SignedEvent = "SignedEvent"
)

type Producer interface {
	PushToQueue(topic string, message []byte) error
}

type Consumer interface {
	Register(ctx context.Context, topic string) chan string
}

type Queue interface {
	Producer
	Consumer
}
