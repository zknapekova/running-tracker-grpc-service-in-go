# Running tracker (gRPC service in Go)

This repository includes gRPC service application that interacts with MongoDB database. 

The fifth prototype includes AddActivities, GetActivities, AddTrainers, GetTrainers, UpdateTrainers, DeleteTrainers API endpoints. It also uses TLS to ensure secure communication between the client and server. The gRPC server is fully dockerized and you can start both server and MongoDB by running 
```
docker compose up grpcserver --build
```
For running tests, execute
```
just unit_tests
```
or 
```
just e2e_tests
```
