package handlers

import (
	"go.mongodb.org/mongo-driver/bson"
	pb "grpcserver/proto/generated_files"
	"reflect"
	"testing"
)

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
