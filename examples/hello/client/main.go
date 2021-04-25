package main

import (
	"context"

	"github.com/saromanov/rabbitmq-rpc/internal/publish"
	"github.com/saromanov/rabbitmq-rpc/examples/hello/proto"
	"github.com/streadway/amqp"
)

func main() {
	ctx := context.Background()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	p, err := publish.New(ch, nil)
	if err != nil {
		panic(err)
	}
	msg := proto.Request{ID: 1}
	res, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	p.Do(ctx, "test", "reply", res)
}
