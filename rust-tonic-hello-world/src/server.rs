mod generated_code {
    tonic::include_proto!("pkg_hello"); 
}

use generated_code::say_server::{Say, SayServer};
use generated_code::{SayRequest, SayResponse};
use tonic::{transport::Server, Request, Response, Status};

#[derive(Debug, Default)]
pub struct MySayer {}

#[tonic::async_trait]
impl Say for MySayer {
    async fn send(&self, request: Request<SayRequest>) -> Result<Response<SayResponse>, Status> {
        println!("Got a request: {:?}", request);

        let reply = generated_code::SayResponse {
            message: format!("Hello {}!", request.into_inner().name).into(),
        };

        Ok(tonic::Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("Hello, world!");

    let addr = "[::1]:50001".parse()?;
    let mysayer = MySayer::default();

    Server::builder()
        .add_service(SayServer::new(mysayer))
        .serve(addr)
        .await?;

    Ok(())
}
