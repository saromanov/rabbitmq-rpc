package rpc

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Handler defines method for handling of request
type Handler func([]byte)([]byte, error)

// Server defines interface for server
type Server interface {
	Close()
	Serve(context.Context, Handler)
}

type Receive struct {
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
	done    chan bool
}

// New creates a new receiver
func New(channel *amqp.Channel, queue string) (Server, error) {
	msgs, err := channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(err, "unable to subscribe to queue")
	}

	return &Receive{channel: channel, msgs: msgs, done: make(chan bool)}, nil
}

// Close provides closing and shutdown of the server
func (srv *Receive) Close() {
	if srv == nil {
		return
	}

	srv.done <- true
}

// Serve provides receiving of the messages
func (srv *Receive) Serve(ctx context.Context, f Handler) {
	finish := false
	for !finish {
		select {
		case msg := <-srv.msgs:
			go srv.callHandler(f, msg)
		case <-ctx.Done():
			return
		case <-srv.done:
			finish = true
		}
	}
}

func (srv *Receive) callHandler(f Handler, msg amqp.Delivery) {

	respData, err := f(msg.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed marshall responce: %v", err))
	}

	err = srv.channel.Publish(
		"",          
		msg.ReplyTo, 
		false,      
		false, 
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: msg.CorrelationId,
			Body:          respData,
		})

	if err != nil {
		panic(fmt.Sprintf("Failed to publish a message: %v", err))
	}

	msg.Ack(false)
}
