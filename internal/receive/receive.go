package rpc

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type CallHandler func(funcID int32, args []byte) ([]byte, error)

type Server interface {
	Close()
	Serve(handler CallHandler)
}

type Receive struct {
	channel *amqp.Channel
	delivery
}

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

	return &serverImpl{conn: conn, channel: ch, msgs: msgs, done: make(chan bool)}, nil
}

type serverImpl struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
	done    chan bool
}

func (srv *serverImpl) Close() {
	if srv == nil {
		return
	}

	srv.done <- true

	if srv.channel != nil {
		srv.channel.Close()
	}

	if srv.conn != nil {
		srv.conn.Close()
	}
}

func (srv *serverImpl) Serve(handler CallHandler) {
	finish := false
	for !finish {
		select {
		case msg := <-srv.msgs:
			go srv.callHandler(handler, msg)
		case <-srv.done:
			finish = true
		}
	}
}

func (srv *serverImpl) callHandler(handler CallHandler, msg amqp.Delivery) {
	var req Request
	err := req.Unmarshal(msg.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed unmarshal request: %v", err))
	}

	var resp Response
	data, err := handler(req.FuncID, req.Body)
	if err != nil {
		resp.IsSuccess = false
		resp.ErrText = err.Error()
	} else {
		resp.IsSuccess = true
		resp.Body = data
	}

	respData, err := resp.Marshal()
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
