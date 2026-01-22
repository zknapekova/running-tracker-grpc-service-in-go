package handlers

import (
	"context"
	"fmt"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	mongodb "grpcserver/mongo_db"
	pb "grpcserver/proto/generated_files"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddTrainers(ctx context.Context, req *pb.AddTrainersRequest) (*pb.AddTrainersResponse, error) {
	request_trainers := req.GetTrainers()

	//validation
	if err := validateAddTrainersRequest(request_trainers); err != nil {
		return nil, err
	}

	//add trainers to DB
	addedTrainers, err := mongodb.AddTrainersToDB(ctx, request_trainers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	utils.Logger.Info("Added trainers", zap.Any("addedTrainers", addedTrainers))

	//extract ids
	ids := make([]string, 0, len(addedTrainers))
	for _, t := range addedTrainers {
		ids = append(ids, t.Id)
	}
	utils.Logger.Info("Extracted ids", zap.Any("ids", ids))

	return &pb.AddTrainersResponse{
		Message: "Trainers were added to database",
		Ids:     ids,
	}, nil
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

	if err := validateUpdateTrainersRequest(req.Trainers); err != nil {
		return nil, err
	}

	updatedTrainers, err := mongodb.UpdateTrainersInDB(ctx, req.Trainers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ids := make([]string, 0, len(updatedTrainers))
	for _, t := range updatedTrainers {
		ids = append(ids, t.Id)
	}
	utils.Logger.Info("Extracted ids", zap.Any("ids", ids))

	return &pb.UpdateTrainersResponse{
		Ids: ids,
	}, nil
}

func (s *Server) DeleteTrainers(ctx context.Context, req *pb.DeleteTrainersRequest) (*pb.DeleteTrainersResponse, error) {
	ids := req.GetIds()

	if len(ids) == 0 {
		return nil, utils.ErrorHandler(nil, "No trainer IDs provided")
	}

	client := mongodb.MongoClient

	objectIds := make([]primitive.ObjectID, len(ids))
	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintf("incorrect ids: %v", id))
		}
		objectIds = append(objectIds, objectId)
	}

	coll := client.Database("data").Collection("trainers")
	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}

	var foundIds []bson.M
	err = cursor.All(ctx, &foundIds)
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}

	len_foundIds := len(foundIds)

	utils.Logger.Info("Number of found ids", zap.Int("ids", len_foundIds))
	if len_foundIds == 0 {
		return nil, utils.ErrorHandler(err, "No trainers to delete were found in DB")
	}

	result, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}

	utils.Logger.Info("Number of deleted docs", zap.Any("deletedCount", result.DeletedCount))
	if result.DeletedCount == 0 {
		return nil, utils.ErrorHandler(err, fmt.Sprintf("DatabaseError: %d trainers found, but no trainers were deleted", len_foundIds))
	}

	deletedIds := make([]string, 0, len_foundIds)
	for _, found_id := range foundIds {
		if id, ok := found_id["_id"].(primitive.ObjectID); ok {
			deletedIds = append(deletedIds, id.Hex())
		}
	}

	return &pb.DeleteTrainersResponse{
		Message: fmt.Sprintf("%d trainer(s) successfully deleted", len_foundIds),
		Ids:     deletedIds,
	}, nil
}
