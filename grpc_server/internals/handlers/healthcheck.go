package handlers

import (
	"context"
	"grpcserver/internals/utils"
	mongodb "grpcserver/mongo_db"
	pb "grpcserver/proto/generated_files"

	"go.uber.org/zap"
)

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {

	status := pb.HealthCheckResponse_SERVING

	err := mongodb.MongoClient.Ping(ctx, nil)
	if err != nil {
		status = pb.HealthCheckResponse_NOT_SERVING
		utils.Logger.Error("MongoDB check failed", zap.Error(err))
	}

	return &pb.HealthCheckResponse{
		Status: status,
	}, nil
}
