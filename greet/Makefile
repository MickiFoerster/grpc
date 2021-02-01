all: rust/src/greet_grpc.rs rust/src/greet.rs

rust/src/greet_grpc.rs : greet.proto /home/micki/.cargo/bin/grpc_rust_plugin
	protoc --grpc_out=rust/src/ --plugin=protoc-gen-grpc=/home/micki/.cargo/bin/grpc_rust_plugin  greet.proto 
rust/src/greet.rs : greet.proto /home/micki/.cargo/bin/protoc-gen-rust
	protoc --rust_out=rust/src/ greet.proto 

/home/micki/.cargo/bin/grpc_rust_plugin:
	cargo install grpcio-compiler
/home/micki/.cargo/bin/protoc-gen-rust:
	cargo install protobuf-codegen

.PHONY: all
