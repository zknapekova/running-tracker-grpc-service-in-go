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

func validateAddActivitiesRequest(request_activities []*pb.Activity) error {
	if len(request_activities) == 0 {
		return status.Error(codes.InvalidArgument, "No activities provided")
	}

	for _, activity := range request_activities {
		if activity.Id != "" {
			return status.Error(codes.InvalidArgument, "Request contains activity with predefined ID")
		}
		// check required fields
		if activity.Name == "" {
			return status.Error(codes.InvalidArgument, "Name field is missing")
		}
		if activity.Duration == 0 {
			return status.Error(codes.InvalidArgument, "Duration field is missing")
		}
		if activity.Date == "" {
			return status.Error(codes.InvalidArgument, "Date field is missing")
		}
	}
	return nil
}
