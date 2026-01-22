package main

import (
	"context"
	"grpcserver/internals/handlers"
	"grpcserver/internals/interceptors"
	"grpcserver/internals/utils"
	mongodb "grpcserver/mongo_db"
	pb "grpcserver/proto/generated_files"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	logger, err := utils.InitLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	utils.Logger = logger

	err = godotenv.Load()
	if err != nil {
		utils.Logger.Fatal("Error loading .env file", zap.Error(err))
	}

	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if certPath == "" || keyPath == "" {
		utils.Logger.Fatal("CERT_PATH or KEY_PATH not set")
	}

	// start TCP listener as TCP is inherently streamed-oriented, establishes connection before data transfer
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		utils.Logger.Fatal("SERVER_PORT not set")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		utils.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	client, err := mongodb.CreateMongoClient()
	if err != nil {
		utils.Logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongodb.DisconnectMongoClient(client, context.Background())

	// initialize gRPC server instance
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		utils.Logger.Fatal("Failed to load credentials", zap.Error(err))
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors.OUAuthentification),
		grpc.Creds(creds),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTrainersServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterActivitiesServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterHealthCheckServiceServer(grpcServer, &handlers.Server{})

	utils.Logger.Info("Server is running", zap.String("port", port))

	// start the server
	err = grpcServer.Serve(listener)
	if err != nil {
		utils.Logger.Fatal("Failed to start the server", zap.Error(err))
	}
}
