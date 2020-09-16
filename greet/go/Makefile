PROTO=$(shell ls ../*.proto)
CLIENT=client/main
SERVER=server/main

all : $(CLIENT) $(SERVER)
	@echo "Finished with build of client ($(CLIENT)) and server ($(SERVER))"

greet.pb.go : $(PROTO)
	protoc -I.. $$(basename $(PROTO)) --go_out=plugins=grpc:.

$(CLIENT) : client/main.go greet.pb.go
	go build -o $@ $<
$(SERVER) : server/main.go greet.pb.go
	go build -o $@ $<

clean : 
	rm -f $(SERVER) $(CLIENT) greet.pb.go

.PHONY: all clean
