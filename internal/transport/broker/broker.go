package broker

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"hezzl/internal/core"
	"log"
)

type LogsService interface {
	ProcessBatch(batch []core.Item) error
}

type MessageBroker interface {
	Sub(subj string, handler nats.MsgHandler) (*nats.Subscription, error)
}

type Broker struct {
	MessageBroker
	LogsService
	Logger *log.Logger
	queue  []core.Item
}

func NewBroker(subj string, m MessageBroker, l LogsService, log *log.Logger) (*Broker, error) {
	b := &Broker{
		MessageBroker: m,
		LogsService:   l,
		Logger:        log,
		queue:         make([]core.Item, 0, 10),
	}

	if _, err := m.Sub(subj, b.updateHandler); err != nil {
		return nil, fmt.Errorf("error while trying to subscribe: %w", err)
	}

	return b, nil
}

func (b *Broker) updateHandler(msg *nats.Msg) {
	if len(b.queue) == 10 {
		if err := b.LogsService.ProcessBatch(b.queue); err != nil {
			b.Logger.Printf("error: %v", err)
			return
		}
		b.queue = b.queue[:0]
	}

	var item core.Item
	if err := json.Unmarshal(msg.Data, &item); err != nil {
		b.Logger.Printf("error: %v", err)
		return
	}
	b.queue = append(b.queue, item)
}
