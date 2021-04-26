export PATH:=$(GOPATH)/bin:$(PATH)

proto-client : 
	protoc -I=. -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf --gofast_out=. ./proto/*.proto