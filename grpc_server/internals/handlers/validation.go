package handlers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpcserver/proto/generated_files"
)

func validateAddTrainersRequest(request_trainers []*pb.Trainer) error {
	if len(request_trainers) == 0 {
		return status.Error(codes.InvalidArgument, "No trainers provided")
	}

	for _, trainer := range request_trainers {
		if trainer.Id != "" {
			return status.Error(codes.InvalidArgument, "request contains trainer with predefined ID")
		}
		if trainer.Brand == "" {
			return status.Error(codes.InvalidArgument, "brand field is missing")
		}
		if trainer.Model == "" {
			return status.Error(codes.InvalidArgument, "model field is missing")
		}
	}
	return nil
}
