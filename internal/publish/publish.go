package publish

import (
	"fmt"

	"github.com/saromanov/rabbitmq-rpc/internal/tools"
	"github.com/streadway/amqp"
	"github.com/pkg/errors"
)
type Publish struct {
	channel *amqp.Channel
}

// New creates new publisher
func New(channel *amqp.Channel) (*Publish, error){
	if channel == nil {
		return nil, errors.New("channel is not defined")
	}
	return &Publish{
		channel: channel,
	}, nil
}

// Do provides sending of the message
func (p *Publish) Do(queue string, data []byte) error {
	err := p.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: tools.GenerateUUID(),
			ReplyTo:       "test-reply",
			Body:          data,
			Expiration:    "1",
		})
	if err != nil {
		return errors.Wrap(err, "unable to send message")
	}
	return nil
}

