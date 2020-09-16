PROTO=$(shell ls *.proto)
PASSWORD_CA_KEY=secretpassword
PASSWORD_SERVER_KEY=secretpassword
SERVER_CN=localhost

.PRECIOUS: %.pb.go
%.pb.go : %.proto
	protoc --go_out=plugins=grpc:. $<

all : grpc_client grpc_server

client/client.go : greet-tls.pb.go ca.crt
grpc_client : client/client.go
	go build -o $@ $<

server/server.go : greet-tls.pb.go server.pem server.crt
grpc_server : server/server.go
	go build -o $@ $<

ca.key :
	openssl genrsa -passout pass:$(PASSWORD_CA_KEY) -des3 -out $@ 2048

ca.crt : ca.key
	openssl req -passin pass:$(PASSWORD_CA_KEY) -new -x509 -days 365 -key $< -out $@ -subj "/CN=$(SERVER_CN)"

server.key :
	openssl genrsa -passout pass:$(PASSWORD_SERVER_KEY) -des3 -out $@ 2048

server.csr : server.key
	openssl req -passin pass:$(PASSWORD_SERVER_KEY) -new -key $< -out $@ -subj "/CN=$(SERVER_CN)"

server.crt : server.csr ca.crt ca.key
	@echo "server signing request: $(word 1,$^)"
	@echo "CA certificate: $(word 2,$^)"
	@echo "CA private key: $(word 3,$^)"
	openssl x509 -req -passin pass:$(PASSWORD_CA_KEY) -days 365 -in $(word 1, $^) -CA $(word 2, $^) -CAkey $(word 3, $^) -set_serial 01 -out $@


server.pem : server.key
	openssl pkcs8 -topk8 -nocrypt -passin pass:$(PASSWORD_SERVER_KEY) -in $< -out $@
