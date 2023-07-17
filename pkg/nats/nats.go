package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"hezzl/config"
)

type NatsBroker struct {
	conn *nats.Conn
}

func NewNatsBroker(cfg config.NatsConfig) (*NatsBroker, error) {
	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("error while trying to connect to nats server: %w", err)
	}

	return &NatsBroker{
		conn: nc,
	}, nil
}

func (n NatsBroker) Publish(subj string, data []byte) error {
	return n.conn.Publish(subj, data)
}

func (n NatsBroker) Sub(subj string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return n.conn.Subscribe(subj, handler)
}
