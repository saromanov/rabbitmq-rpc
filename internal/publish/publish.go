package publish

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/saromanov/rabbitmq-rpc/internal/models"
	"github.com/saromanov/rabbitmq-rpc/internal/tools"
	"github.com/streadway/amqp"
)
type Publish struct {
	mu sync.Mutex
	unmarshaller Unmarshaller
	channel *amqp.Channel
	calls map[string]*models.Call
}

// New creates new publisher
func New(channel *amqp.Channel, um Unmarshaller) (*Publish, error){
	if channel == nil {
		return nil, errors.New("channel is not defined")
	}
	return &Publish{
		mu: sync.Mutex{},
		unmarshaller: um,
		channel: channel,
		calls: make(map[string]*models.Call),
	}, nil
}

// Do provides sending of the message
func (p *Publish) Do(ctx context.Context, queue, replyQueue string, data []byte) ([]byte, error) {
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
		return nil, errors.Wrap(err, "unable to send message")
	}
	return p.handleCall(ctx, corrID)
}

func (p *Publish) handleCall(ctx context.Context, corrID string) ([]byte, error) {
	call := &models.Call{Done: make(chan bool)}
	p.mu.Lock()
	p.calls[corrID] = call
	p.mu.Unlock()

	var resp []byte
	select {
	case <- call.Done:
		if err := p.unmarshaller.Do(call.Data, nil); err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal data")
		}
		return nil, nil
	case <- ctx.Done():
		return nil, errors.New("unable to get call")
	}
	p.mu.Lock()
	delete(p.calls, corrID)
	p.mu.Unlock()
	return resp, nil
}

