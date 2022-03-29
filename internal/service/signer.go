package service

import (
	"context"
	"encoding/base64"
	"github.com/enverbisevac/xdefi/internal/config"
	"github.com/enverbisevac/xdefi/internal/queue"
	"log"
)

type Signer struct {
	queue queue.Queue
}

func NewSigner(q queue.Queue) *Signer {
	return &Signer{
		queue: q,
	}
}

func (s *Signer) Start(ctx context.Context) {
	signChan := s.queue.Register(ctx, queue.SignEvent)
	go func() {
		for data := range signChan {
			log.Println(data)
			err := s.queue.PushToQueue(queue.SignedEvent, []byte(SignEvent(data)))
			if err != nil {
				log.Printf("error while sending data to the queue %v", err)
			}
		}
	}()
}

func SignEvent(input string) string {
	key := config.GetXORKey()
	bytes := make([]byte, 0)
	for i := 0; i < len(input); i++ {
		bytes = append(bytes, input[i]^key[i%len(key)])
	}

	return base64.StdEncoding.EncodeToString(bytes)
}
