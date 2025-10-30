package handlers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
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

func (s *Server) UpdateTrainers(ctx context.Context, req *pb.UpdateTrainersRequest) (*pb.UpdateTrainersResponse, error) {
	client, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	var updatedTrainers []*pb.Trainer
	for _, trainer := range req.Trainers {
		modelTrainer := mongodb.MapPbTrainersToModelTrainers(trainer)

		objId, err := primitive.ObjectIDFromHex(trainer.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid Id")
		}

		modelDoc, err := bson.Marshal(modelTrainer)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal error")
		}

		var updateDoc bson.M
		err = bson.Unmarshal(modelDoc, &updateDoc)
		if err != nil {
			return nil, utils.ErrorHandler(err, "internal error")
		}

		delete(updateDoc, "_id")

		_, err = client.Database("main").Collection("trainers").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintln("error updating teacher id:", trainer.Id))
		}

		updatedTrainer := mongodb.MapModelTrainersToPb(modelTrainer)

		updatedTrainers = append(updatedTrainers, updatedTrainer)

	}
	ids := make([]string, 0, len(updatedTrainers))
	for _, t := range updatedTrainers {
		ids = append(ids, t.Id)
	}
	fmt.Println(ids)

	return &pb.UpdateTrainersResponse{
		Ids: ids,
	}, nil

}
