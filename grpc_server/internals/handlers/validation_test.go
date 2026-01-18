package handlers

import (
	pb "grpcserver/proto/generated_files"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateAddTrainersRequest(t *testing.T) {
	tests := []struct {
		name            string
		request         []*pb.Trainer
		expectedMessage string
	}{
		{
			name: "missing brand",
			request: []*pb.Trainer{
				{
					Model:        "test model",
					PurchaseDate: "1999-12-04",
				},
			},
		},
		{
			name: "missing model",
			request: []*pb.Trainer{
				{
					Brand:        "test brand",
					PurchaseDate: "1999-12-04",
				},
			},
		},
		{
			name: "predefined id",
			request: []*pb.Trainer{
				{
					Brand:        "test brand",
					PurchaseDate: "1999-12-04",
					Id:           "1000",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateAddTrainersRequest(tt.request)

			if res == nil {
				t.Fatalf("expected error but got nil")
			}

			st, ok := status.FromError(res)
			if !ok {
				t.Fatalf("expected gRPC status error, got %v", res)
			}

			if st.Code() != codes.InvalidArgument {
				t.Fatalf("expected code InvalidArgument, got %v", st.Code())
			}

			if tt.expectedMessage != "" && st.Message() != tt.expectedMessage {
				t.Fatalf("expected message %q, got %q", tt.expectedMessage, st.Message())
			}

		})
	}
}

func TestValidateUpdateTrainersRequest(t *testing.T) {
	tests := []struct {
		name            string
		request         []*pb.Trainer
		expectedMessage string
	}{
		{
			name:            "no data",
			request:         []*pb.Trainer{},
			expectedMessage: "No trainers provided",
		},
		{
			name: "empty id",
			request: []*pb.Trainer{
				{
					Id: "",
				},
			},
			expectedMessage: "No id specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateUpdateTrainersRequest(tt.request)

			if res == nil {
				t.Fatalf("expected error but got nil")
			}

			st, ok := status.FromError(res)
			if !ok {
				t.Fatalf("expected gRPC status error, got %v", res)
			}

			if st.Code() != codes.InvalidArgument {
				t.Fatalf("expected code %v, got %v", codes.InvalidArgument, st.Code())
			}

			if tt.expectedMessage != "" && st.Message() != tt.expectedMessage {
				t.Fatalf("expected message %q, got %q", tt.expectedMessage, st.Message())
			}
		})
	}
}

func TestValidateAddActivitiesRequest(t *testing.T) {
	tests := []struct {
		name            string
		request         []*pb.Activity
		expectedError   bool
		expectedMessage string
	}{
		{
			name:            "no data",
			request:         []*pb.Activity{},
			expectedMessage: "No activities provided",
		},
		{
			name: "predefined id",
			request: []*pb.Activity{
				{
					Id:       "10000",
					Duration: 122,
					Date:     "01-01-2026",
					Name:     "swimming",
				},
			},
			expectedMessage: "Request contains activity with predefined ID",
		},
		{
			name: "missing duration",
			request: []*pb.Activity{
				{
					Date: "01-01-2026",
					Name: "swimming",
				},
			},
			expectedMessage: "Duration field is missing",
		},
		{
			name: "missing name",
			request: []*pb.Activity{
				{
					Duration: 120,
					Date:     "01-01-2026",
				},
			},
			expectedMessage: "Name field is missing",
		},
		{
			name: "missing date",
			request: []*pb.Activity{
				{
					Duration: 120,
					Name:     "hiking",
				},
			},
			expectedMessage: "Date field is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateAddActivitiesRequest(tt.request)

			if res == nil {
				t.Fatalf("returned nil instead of InvalidArgument")
			}

			st, ok := status.FromError(res)
			if !ok {
				t.Fatalf("expected gRPC status error, got %v", res)
			}

			if st.Code() != codes.InvalidArgument {
				t.Fatalf("expected InvalidArgument, got %v", st.Code())
			}

			if tt.expectedMessage != "" && st.Message() != tt.expectedMessage {
				t.Fatalf("expected message %q, got %q", tt.expectedMessage, st.Message())
			}
		})
	}
}
