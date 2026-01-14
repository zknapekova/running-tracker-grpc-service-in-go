package handlers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpcserver/proto/generated_files"
	"testing"
)

func TestValidateAddTrainersRequest_MissingBrand(t *testing.T) {
	// Test that validateAddTrainersRequest function returns Invalid argument error when brand parameter is not specified

	req := []*pb.Trainer{
		{
			Model:        "test model",
			PurchaseDate: "1999-12-04",
		},
	}

	res := validateAddTrainersRequest(req)

	if res == nil {
		t.Fatalf("validateAddTrainersRequest returned nil instead of InvalidArgument")
	}

	st, ok := status.FromError(res)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", res)
	}

	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestValidateAddTrainersRequest_MissingModel(t *testing.T) {
	// Ensure validateAddTrainersRequest returns Invalid argument error code when model parameter is missing

	req := []*pb.Trainer{
		{
			Brand:        "test brand",
			PurchaseDate: "1999-12-04",
		},
	}

	res := validateAddTrainersRequest(req)

	if res == nil {
		t.Fatalf("validateAddTrainersRequest returned nil instead of InvalidArgument")
	}

	st, ok := status.FromError(res)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", res)
	}

	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", st.Code())
	}
}
