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
	"reflect"
	"strings"
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
	filter, err := buildFilterForTrainers(req.Trainers)
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

func buildFilterForTrainers(train_obj *pb.Trainer) (bson.M, error) {
	filter := bson.M{}

	if train_obj == nil {
		return filter, nil
	}

	var modelTrainers models.Trainers
	modelVal := reflect.ValueOf(&modelTrainers).Elem()
	modelType := modelVal.Type()

	reqVal := reflect.ValueOf(train_obj).Elem()
	reqType := reqVal.Type()

	for i := 0; i < reqVal.NumField(); i++ {
		fieldVal := reqVal.Field(i)
		fieldName := reqType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(fieldVal)
			}
		}
	}
	for i := 0; i < modelVal.NumField(); i++ {
		fieldVal := modelVal.Field(i)

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")
			if bsonTag == "_id" {
				objId, err := primitive.ObjectIDFromHex(train_obj.Id)
				if err != nil {
					return nil, utils.ErrorHandler(err, "Internal Error")
				}
				filter[bsonTag] = objId
			} else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}
	fmt.Println(filter)
	return filter, nil
}

func buildSortOptions(sortFields []*pb.SortField) bson.D {
	var sortOptions bson.D
	for _, sortField := range sortFields {
		order := 1
		if sortField.GetOrder() == pb.Order_DESC {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{Key: sortField.Field, Value: order})
	}
	fmt.Println("Sort options", sortOptions)
	return sortOptions
}
