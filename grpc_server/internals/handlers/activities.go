package handlers

import (
	"context"
	"grpcserver/internals/utils"
	mongodb "grpcserver/mongo_db"
	pb "grpcserver/proto/generated_files"

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
	utils.InfoLogger.Println(addedActivities)

	//extract ids
	ids := make([]string, 0, len(addedActivities))
	for _, t := range addedActivities {
		ids = append(ids, t.Id)
	}
	utils.InfoLogger.Println(ids)

	return &pb.AddActivitiesResponse{
		Message: "Activities were added to the database",
		Ids:     ids,
	}, nil
}
