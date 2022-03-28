package queue

type Producer interface {
	PushToQueue(message []byte) error
}
