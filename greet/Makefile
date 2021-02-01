all: greet_grpc.rs greet.rs


greet_grpc.rs : greet.proto
	protoc --grpc_out=rust/src/ --plugin=protoc-gen-grpc=/home/micki/.cargo/bin/grpc_rust_plugin  greet.proto 
greet.rs : greet.proto
	protoc --rust_out=rust/src/ greet.proto 
