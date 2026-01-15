package interceptors

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"strings"
)

func OUAuthentification(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	token := strings.TrimSpace(os.Getenv("OAUTH_TOKEN"))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata unavailable")
	}
	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token unavailable")
	}
	ctx_token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if ctx_token != token {
		return nil, status.Errorf(codes.PermissionDenied, "Incorrect token: Permission denied")
	}
	return handler(ctx, req)
}
