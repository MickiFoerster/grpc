#include <grpcpp/impl/codegen/client_context.h>
#include <iostream>
#include <memory>
#include <string>

#include <grpcpp/grpcpp.h>

#include "greet.grpc.pb.h"

//using grpc::Channel;
//using grpc::ClientContext;
//using grpc::Status;
//using greet::GreetService;

class GreetClient
{
public:
  GreetClient(std::shared_ptr<grpc::Channel> channel)
      : stub_(greet::GreetService::NewStub(channel)) {}

  std::string Greet(std::string firstname, std::string backname)
  {
    auto greeting = greet::Greeting().New();
    greeting->set_first_name("John");
    greeting->set_last_name("Doo");

    greet::GreetRequest request;
    request.set_allocated_greeting(greeting);

    greet::GreetResponse response;
    grpc::ClientContext context;
    grpc::Status status;

    status = stub_->Greet(&context, request, &response);
    if (status.ok())
    {
      return response.result();
    }
    else
    {
      std::cout << status.error_code()
                << ": " << status.error_message()
                << std::endl;
      return "RPC failed";
    }
  }

private:
  std::unique_ptr<greet::GreetService::Stub> stub_;
};

int main()
{
  auto chan = grpc::CreateChannel("localhost:50051",
                                  grpc::InsecureChannelCredentials());

  GreetClient greetclient(chan);
  std::string reply = greetclient.Greet("John", "Doo");
  std::cout << "Greet RPC returned: " << reply << std::endl;

  return 0;
}
