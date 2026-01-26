package handlers

import (
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	pb "grpcserver/proto/generated_files"
	"os"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	utils.Logger = zap.NewNop()

	os.Exit(m.Run())
}

func TestBuildSortOptions_Empty(t *testing.T) {
	// Check that buildSortOptions returns an empty document if no sort option is specified

	input := []*pb.SortField{}

	result := buildSortOptions(input)
	if len(result) != 0 {
		t.Fatalf("Expected empty doc, got %v", result)
	}
}

func TestBuildSortOptions_Desc(t *testing.T) {
	// Check that buildSortOptions returns the correct document when descending order is requested

	field_name := "test"
	input := []*pb.SortField{
		{
			Field: field_name,
			Order: pb.Order_DESC,
		},
	}

	result := buildSortOptions(input)
	expect := bson.D{{Key: field_name, Value: -1}}
	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("Expected %v, got %v", expect, result)
	}
}

func TestBuildSortOptions_Asc(t *testing.T) {
	// Verify that buildSortOptions returns the correct document when ascending order is requested

	field_name := "test"
	input := []*pb.SortField{
		{
			Field: field_name,
			Order: pb.Order_ASC,
		},
	}

	result := buildSortOptions(input)
	expect := bson.D{{Key: field_name, Value: 1}}
	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("Expected %v, got %v", expect, result)
	}
}

func TestBuildSortOptions_DescAsc(t *testing.T) {
	// Verify that buildSortOptions returns the correct document when both descending and ascending order are requested

	field_name_desc := "test_desc"
	field_name_asc := "test_asc"

	input := []*pb.SortField{
		{
			Field: field_name_desc,
			Order: pb.Order_DESC,
		},
		{
			Field: field_name_asc,
			Order: pb.Order_ASC,
		},
	}

	result := buildSortOptions(input)
	expect := bson.D{{Key: field_name_desc, Value: -1}, {Key: field_name_asc, Value: 1}}
	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("Expected %v, got %v", expect, result)
	}
}

func TestBuildFilter_EmptyFilter(t *testing.T) {
	// Ensure that buildFilter function returns empty filter map when no filters were requested

	input := &pb.Trainer{}

	filter, err := buildFilter(input, &models.Trainers{})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(filter) != 0 {
		t.Fatalf("Expected empty map, got %v", filter)
	}
}

func TestBuildFilter_NilInput(t *testing.T) {
	// Ensure that buildFilter function returns empty filter map when no filters were requested

	filter, err := buildFilter(nil, &models.Trainers{})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(filter) != 0 {
		t.Fatalf("Expected empty map, got %v", filter)
	}
}

func TestBuildFilter_AllFilters(t *testing.T) {
	// Ensure that buildFilter function returns the correct filter map when all fields are provided

	input := &pb.Trainer{
		Brand:            "test_brand",
		ExpectedLifespan: 700,
		Model:            "test_model",
		PurchaseDate:     "2026-01-01",
		SurfaceType:      pb.SurfaceType_ROAD,
		Status:           pb.TrainerStatus_ACTIVE,
	}

	filter, err := buildFilter(input, &models.Trainers{})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	expect := bson.M{
		"brand":             "test_brand",
		"expected_lifespan": int64(700),
		"model":             "test_model",
		"purchase_date":     "2026-01-01",
		"surface_type":      pb.SurfaceType_ROAD,
		"status":            pb.TrainerStatus_ACTIVE,
	}

	if !reflect.DeepEqual(filter, expect) {
		t.Fatalf("Expected %v, got %v", expect, filter)
	}
}

func TestBuildFilter_EmptyFieldValue(t *testing.T) {
	// Check if buildFilter function omits fileds with empty values

	input := &pb.Trainer{
		Brand:            "",
		ExpectedLifespan: 700,
		Model:            "",
		PurchaseDate:     "2026-01-01",
		SurfaceType:      pb.SurfaceType_ROAD,
		Status:           pb.TrainerStatus_ACTIVE,
	}

	filter, err := buildFilter(input, &models.Trainers{})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	expect := bson.M{
		"expected_lifespan": int64(700),
		"purchase_date":     "2026-01-01",
		"surface_type":      pb.SurfaceType_ROAD,
		"status":            pb.TrainerStatus_ACTIVE,
	}

	if !reflect.DeepEqual(filter, expect) {
		t.Fatalf("Expected %v, got %v", expect, filter)
	}
}
