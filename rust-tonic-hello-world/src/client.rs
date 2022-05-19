mod generated_code {
    tonic::include_proto!("pkg_hello");
}

use generated_code::say_client::SayClient;
use generated_code::SayRequest;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = SayClient::connect("http://[::1]:50001").await?;

    let request = tonic::Request::new(SayRequest {
        name: "Micki".into(),
    });

    let response = client.send(request).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
