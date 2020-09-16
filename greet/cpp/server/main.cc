#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include <grpcpp/grpcpp.h>
#include <grpcpp/health_check_service_interface.h>

#include "greet.grpc.pb.h"

class GreetServiceImpl final : public greet::GreetService::Service {
  grpc::Status Greet(grpc::ServerContext *context,
                     const greet::GreetRequest *request,
                     greet::GreetResponse *response) override {
    std::string prefix("Hello ");
    std::string reply(prefix + request->greeting().first_name() + " " +
                      request->greeting().last_name());
    std::cout << "Send '" << reply << "' to the client." << std::endl;
    response->set_result(reply);
    return grpc::Status::OK;
  }
};

int main() {
  std::string server_address("0.0.0.0:50051");
  GreetServiceImpl service;

  grpc::EnableDefaultHealthCheckService(true);
  grpc::reflection::InitProtoReflectionServerBuilderPlugin();
  grpc::ServerBuilder builder;

  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);
  std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
  std::cout << "Server listening on " << server_address << std::endl;

  // Wait for the server to shutdown. Note that some other thread must be
  // responsible for shutting down the server for this call to ever return.
  server->Wait();

  return 0;
}

