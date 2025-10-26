package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpcserver/internals/models"
	mongodb "grpcserver/mongo_db"
	pb "grpcserver/proto/generated_files"
)

func (s *Server) AddTrainers(ctx context.Context, req *pb.AddTrainersRequest) (*pb.AddTrainersResponse, error) {
	request_trainers := req.GetTrainers()

	//validation
	if err := validateTrainersRequest(request_trainers); err != nil {
		return nil, err
	}

	//add trainers to DB
	addedTrainers, err := mongodb.AddTrainersToDB(ctx, request_trainers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fmt.Println(addedTrainers)

	//extract ids
	ids := make([]string, 0, len(addedTrainers))
	for _, t := range addedTrainers {
		ids = append(ids, t.Id)
	}
	fmt.Println(ids)

	return &pb.AddTrainersResponse{
		Message: "Trainers were added to database",
		Ids:     ids,
	}, nil
}

func validateTrainersRequest(request_trainers []*pb.Trainer) error {
	if len(request_trainers) == 0 {
		return status.Error(codes.InvalidArgument, "No trainers provided")
	}

	for _, trainer := range request_trainers {
		if trainer.Id != "" {
			return status.Error(codes.InvalidArgument, "request contains trainer with predefined ID")
		}
	}
	return nil
}

func (s *Server) GetTrainers(ctx context.Context, req *pb.GetTrainersRequest) (*pb.GetTrainersResponse, error) {
	filter, err := buildFilter(req.Trainers, &models.Trainers{})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	sortOptions := buildSortOptions(req.GetSortBy())

	trainers, err := mongodb.GetTrainersFromDb(ctx, sortOptions, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetTrainersResponse{Trainers: trainers}, nil
}
