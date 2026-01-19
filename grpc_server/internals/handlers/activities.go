package handlers

import (
	"context"
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
