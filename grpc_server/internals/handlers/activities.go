package handlers

import (
	"context"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	mongodb "grpcserver/mongo_db"
	pb "grpcserver/proto/generated_files"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddActivities(ctx context.Context, req *pb.AddActivitiesRequest) (*pb.AddActivitiesResponse, error) {
	request_activities := req.GetActivities()

	//request validation
	if err := validateAddActivitiesRequest(request_activities); err != nil {
		return nil, err
	}

	//add activities to DB
	addedActivities, err := mongodb.AddActivitiesToDB(ctx, request_activities)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//extract ids
	ids := make([]string, 0, len(addedActivities))
	for _, t := range addedActivities {
		ids = append(ids, t.Id)
	}
	utils.Logger.Info("Extracted ids", zap.Any("ids", ids))

	return &pb.AddActivitiesResponse{
		Message: "Activities were added to the database",
		Ids:     ids,
	}, nil
}

func (s *Server) GetActivities(ctx context.Context, req *pb.GetActivitiesRequest) (*pb.GetActivitiesResponse, error) {
	activity_filter := req.GetActivityFilter()
	if activity_filter == nil {
		return nil, status.Error(codes.InvalidArgument, "No filter specified")
	}

	filter, err := buildFilter(activity_filter, &models.Activities{})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	sortOptions := buildSortOptions(req.GetSortBy())

	activities, err := mongodb.GetActivitiessFromDb(ctx, sortOptions, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetActivitiesResponse{Activities: activities}, nil
}
