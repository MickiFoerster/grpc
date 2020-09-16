PROTO=$(shell ls *.proto)

all : 
	protoc $(PROTO) --go_out=plugins=grpc:.
