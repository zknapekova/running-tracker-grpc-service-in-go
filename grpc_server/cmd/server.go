package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpcserver/internals/handlers"
	"grpcserver/internals/utils"
	"grpcserver/internals/interceptors"
	pb "grpcserver/proto/generated_files"
	"log"
	"net"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if certPath == "" || keyPath == "" {
		log.Fatal("CERT_PATH or KEY_PATH not set")
	}

	// start TCP listener as TCP is inherently streamed-oriented, establishes connection before data transfer
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Fatal("SERVER_PORT not set")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	// initialize gRPC server instance
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		log.Fatal("Failed to load credentials: ", err)
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors.OUAuthentification),
		grpc.Creds(creds),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTrainersServiceServer(grpcServer, &handlers.Server{})

	utils.InfoLogger.Println("Server is running on port", port)

	// start the server
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
