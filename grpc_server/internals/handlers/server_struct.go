package handlers

import pb "grpcserver/proto/generated_files"

type Server struct {
	pb.UnimplementedTrainersServiceServer
}