# Running tracker (gRPC service in Go)

This repository includes gRPC service that receives data from test client and updates MongoDB database. 

The fourth prototype currently includes AddActivities, AddTrainers, GetTrainers, UpdateTrainers, DeleteTrainers gRPC API endpoints. It also uses TLS to ensure secure communication between client and server. The gRPC server is fully dockerized and you can start both server and MongoDB by running 
```
docker compose up grpcserver --build
```
For running tests, navigate to the grpc_server directory and execute
```
just unit_tests
```
or 
```
just e2e_tests
```
