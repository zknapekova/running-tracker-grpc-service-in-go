package main

import  (
		"context"
		"net"
		"log"
		"google.golang.org/grpc"
		pb "grpcserver/proto/generated_files"
)

type server struct {
	pb.UnimplementedTrainersServiceServer
}

func (s *server)AddTrainers(ctx context.Context, in *pb.AddTrainersRequest) (*pb.Response, error) {
	//TODO: handle logic
	return &pb.Response{
		Message: "all good",
		Code: 0,
	}, nil
}


func main () {

	// start TCP listener as TCP is inherently streamed-oriented, establishes connection before data transfer
	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}
	// initialize gRPC server instance
	grpcServer := grpc.NewServer()
	pb.RegisterTrainersServiceServer(grpcServer, &server{})

	log.Println("Server is running on port", port)

	// start the server
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}