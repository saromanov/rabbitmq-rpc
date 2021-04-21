package publish

import (
	"fmt"

	"github.com/saromanov/rabbitmq-rpc/internal/tools"
	"github.com/saromanov/rabbitmq-rpc/internal/models"
	"github.com/streadway/amqp"
	"github.com/pkg/errors"
)
type Publish struct {
	channel *amqp.Channel
	calls map[string]*models.Call
}

// New creates new publisher
func New(channel *amqp.Channel) (*Publish, error){
	if channel == nil {
		return nil, errors.New("channel is not defined")
	}
	return &Publish{
		channel: channel,
		calls: make(map[string]*models.Call),
	}, nil
}

// Do provides sending of the message
func (p *Publish) Do(queue, replyQueue string, data []byte) error {
	corrID := tools.GenerateUUID()
	err := p.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: corrID,
			ReplyTo:       replyQueue,
			Body:          data,
			Expiration:    "1",
		})
	if err != nil {
		return errors.Wrap(err, "unable to send message")
	}
	return p.handleCall(corrID)
}

func (p *Publish) handleCall(corrID string) error {
	p.calls[corrID] = &models.Call{Done: make(chan bool)}
	return nil
}

