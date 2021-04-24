package main 

import (
	"github.com/saromanov/rabbitmq-rpc/internal/publish"
)

func main(){
	p := publish.New()
	p.Do()
}