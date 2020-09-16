PROTO=$(shell ls ../*.proto)
CLIENT=client/main
SERVER=server/main
CXXFLAGS=-ggdb3 -O0 -I. `pkg-config --cflags protobuf grpc++`
LDFLAGS += -L/usr/local/lib `pkg-config --libs protobuf grpc++`\
           -pthread\
           -Wl,--no-as-needed -lgrpc++_reflection -Wl,--as-needed\
           -ldl
OBJFILES=greet.pb.o greet.grpc.pb.o

all : $(CLIENT) $(SERVER)
	@echo "Finished with build of client ($(CLIENT)) and server ($(SERVER))"

greet.grpc.pb.cc : $(PROTO)
	protoc -I.. $$(basename $(PROTO)) --grpc_out=. --plugin=protoc-gen-grpc=`which grpc_cpp_plugin`
greet.pb.cc : $(PROTO)
	protoc -I.. $$(basename $(PROTO)) --cpp_out=.

$(CLIENT) : client/main.cc $(OBJFILES)
	$(CXX) $(CXXFLAGS) -o $@ $< $(OBJFILES) $(LDFLAGS)
$(SERVER) : server/main.cc greet.pb.o greet.grpc.pb.o
	$(CXX) $(CXXFLAGS) -o $@ $< $(OBJFILES) $(LDFLAGS)

clean : 
	rm -f $(SERVER) $(CLIENT) *.pb.h *.pb.cc $(OBJFILES)

.PHONY: all clean
