# Running tracker (gRPC service in Go)

This repository includes gRPC service that receives data from test client and updates MongoDB database. 

The third prototype currently includes add_trainers, get_trainers, update_trainers and delete_trainers endpoints. It also uses TLS to ensure secure communication between client and server. The gRPC server is fully containerized and you can start both server and MongoDB by running 
```
docker compose up grpcserver --build
```
To run the test requests from the client, execute
```
cd grpc_client
go run cmd/client.go
```
