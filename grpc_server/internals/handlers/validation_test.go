package handlers

import (
	pb "grpcserver/proto/generated_files"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestValidateAddTrainersRequest_PredefinedId(t *testing.T) {
	// Check validateAddTrainersRequest returns Invalid argument error when predefined id is included in the request

	req := []*pb.Trainer{
		{
			Brand:        "test brand",
			PurchaseDate: "1999-12-04",
			Id:           "1000",
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

func TestValidateUpdateTrainersRequest_NoData(t *testing.T) {
	// Check validateUpdateTrainersRequest returns Invalid argument error no data is sent

	req := []*pb.Trainer{}

	res := validateUpdateTrainersRequest(req)

	if res == nil {
		t.Fatalf("validateUpdateTrainersRequest returned nil instead of InvalidArgument")
	}

	st, ok := status.FromError(res)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", res)
	}

	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", st.Code())
	}

	expected := "No trainers provided"
	if st.Message() != expected {
		t.Fatalf("expected %s message, got %v", expected, st.Message())
	}
}

func TestValidateUpdateTrainersRequest_EmptyId(t *testing.T) {
	// Test that validateUpdateTrainersRequest returns Invalid argument error when id field is empty

	req := []*pb.Trainer{
		{
			Id: "",
		},
	}

	res := validateUpdateTrainersRequest(req)

	if res == nil {
		t.Fatalf("validateUpdateTrainersRequest returned nil instead of InvalidArgument")
	}

	st, ok := status.FromError(res)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", res)
	}

	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", st.Code())
	}

	expected := "No id specified"
	if st.Message() != expected {
		t.Fatalf("expected %s message, got %v", expected, st.Message())
	}
}
