package handlers

import (
	pb "grpcserver/proto/generated_files"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateAddTrainersRequest(request_trainers []*pb.Trainer) error {
	if len(request_trainers) == 0 {
		return status.Error(codes.InvalidArgument, "No trainers provided")
	}

	for _, trainer := range request_trainers {
		if trainer.Id != "" {
			return status.Error(codes.InvalidArgument, "Request contains trainer with predefined ID")
		}
		if trainer.Brand == "" {
			return status.Error(codes.InvalidArgument, "Brand field is missing")
		}
		if trainer.Model == "" {
			return status.Error(codes.InvalidArgument, "Model field is missing")
		}
	}
	return nil
}

func validateUpdateTrainersRequest(request_trainers []*pb.Trainer) error {

	if len(request_trainers) == 0 {
		return status.Error(codes.InvalidArgument, "No trainers provided")
	}

	for _, trainer := range request_trainers {
		if trainer.Id == "" {
			return status.Error(codes.InvalidArgument, "No id specified")
		}
	}

	return nil

}
