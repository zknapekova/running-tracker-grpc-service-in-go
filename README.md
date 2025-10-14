# Running tracker (gRPC service in Go)

This repository includes gRPC service that receives data from client and updates MongoDB database. 

The first prototype currently includes add_trainers RPC. It also uses TLS to ensure secure communication between client and server. The gRPC server is fully containerized and you can start both server and MongoDB by running `docker compose up --build`.
