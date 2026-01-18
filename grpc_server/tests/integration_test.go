package tests

import (
	"context"
	"grpcserver/client"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mongodb "grpcserver/mongo_db"
	running_trackerpb "grpcserver/proto/generated_files"

	"github.com/joho/godotenv"
)

var (
	new_client   *client.Client
	ctx          context.Context
	mongo_client *mongo.Client
	db           *mongo.Database
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
	}

	certPath := os.Getenv("CERT_PATH")
	token := os.Getenv("OAUTH_TOKEN")

	config := client.Config{
		CertPath:   certPath,
		OAuthToken: token,
	}

	new_client, err = client.CreateServiceClient(config)
	if err != nil {
		log.Fatal("Connection to service failed ", err)
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)

	mongo_client, err = mongodb.CreateMongoClient()
	if err != nil {
		log.Fatal("Connection to database failed ", err)
	}
	db = mongo_client.Database("data")

	code := m.Run()

	new_client.Conn.Close()
	mongodb.DisconnectMongoClient(mongo_client, ctx)
	cancel()
	os.Exit(code)
}

func TestAddTrainers_ok(t *testing.T) {
	// The test sends ADDTrainers request, checks the response and verifies that data were inserted in DB

	add_trainers_request := running_trackerpb.AddTrainersRequest{
		Trainers: []*running_trackerpb.Trainer{
			{
				Brand:            "test_brand1",
				Model:            "test_model1",
				PurchaseDate:     "2024-02-04",
				ExpectedLifespan: 700,
				SurfaceType:      running_trackerpb.SurfaceType_ROAD_TO_TRAIL,
				Status:           running_trackerpb.TrainerStatus_NEW,
			},
			{
				Brand:            "test_brand1",
				Model:            "test_model2",
				PurchaseDate:     "2025-01-01",
				ExpectedLifespan: 700,
				SurfaceType:      running_trackerpb.SurfaceType_TRAIL,
				Status:           running_trackerpb.TrainerStatus_NEW,
			},
		},
	}
	res_add, err := new_client.Trainers.AddTrainers(ctx, &add_trainers_request)
	if err != nil {
		t.Fatal("Could not add ", err)
	}

	log.Println("IDs:", res_add.Ids)

	//db check
	objectIds := make([]primitive.ObjectID, len(res_add.Ids))
	for _, id := range res_add.Ids {
		objectId, _ := primitive.ObjectIDFromHex(id)
		objectIds = append(objectIds, objectId)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	// db clean up
	t.Cleanup(func() {
		_, err_delete := db.Collection("trainers").DeleteMany(ctx, filter)
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete trainers: %v", err)
		}
	})

	cursor, err := db.Collection("trainers").Find(context.Background(), filter)
	if err != nil {
		t.Fatalf("Could not query the DB: %s ", err)
	}
	var foundIds []bson.M
	err = cursor.All(ctx, &foundIds)
	if err != nil {
		t.Fatalf("Decoding failed: %s ", err)
	}

	if len(foundIds) != len(add_trainers_request.Trainers) {
		t.Fatalf("Number of found IDs in DB does not match number of added trainers")
	}

	//response check
	expected_message := "Trainers were added to database"
	if expected_message != res_add.Message {
		t.Fatalf("Expected message %s, got %s", expected_message, res_add.Message)
	}

	if len(res_add.Ids) != len(add_trainers_request.Trainers) {
		t.Fatalf("Number of returned IDs does not match number of added trainers")
	}
}

func TestGetTrainers_SortDesc(t *testing.T) {
	// The test validates that GetTrainers requests returns inserted data and verifies that sorting in descending is perfomed correctly

	// db set up
	test_trainers1 := insertTestTrainer(t, mongo_client, "test_brand1", "model1")
	test_trainers2 := insertTestTrainer(t, mongo_client, "test_brand1", "model2")

	t.Cleanup(func() {
		filter := bson.M{"_id": test_trainers1.ID}
		_, err_delete := db.Collection("trainers").DeleteOne(context.Background(), filter)
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete trainer1: %v", err_delete)
		}

		filter = bson.M{"_id": test_trainers2.ID}
		_, err_delete = db.Collection("trainers").DeleteOne(context.Background(), filter)
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete trainer2: %v", err_delete)
		}
	})

	get_trainers_request := running_trackerpb.GetTrainersRequest{
		Trainers: &running_trackerpb.Trainer{
			Brand: "test_brand1",
		},
		SortBy: []*running_trackerpb.SortField{
			{
				Field: "model",
				Order: running_trackerpb.Order_DESC,
			},
		},
	}

	// response check
	res_get, err := new_client.Trainers.GetTrainers(ctx, &get_trainers_request)
	if err != nil {
		t.Fatal("Could not get ", err)
	}
	log.Println("GET response:", res_get)

	if len(res_get.Trainers) == 0 {
		t.Fatalf("No trainers returned from GetTrainers")
	}

	expectedCount := 2
	if len(res_get.Trainers) != expectedCount {
		t.Fatalf("Expected %d trainers, got %d", expectedCount, len(res_get.Trainers))
	}

	// check sorting
	if res_get.Trainers[0].Model != test_trainers2.Model {
		t.Fatalf("Expected first trainer model to be %s, got '%s'", test_trainers2.Model, res_get.Trainers[0].Model)
	}
	if res_get.Trainers[1].Model != test_trainers1.Model {
		t.Fatalf("Expected second trainer model to be %s, got '%s'", test_trainers1.Model, res_get.Trainers[1].Model)
	}
}

func TestGetTrainers_SortAsc(t *testing.T) {
	// The test validates that GetTrainers requests returns inserted data and verifies that sorting in ascending order is perfomed correctly

	// db set up
	test_trainers1 := insertTestTrainer(t, mongo_client, "test_brand1", "model1")
	test_trainers2 := insertTestTrainer(t, mongo_client, "test_brand1", "model2")

	t.Cleanup(func() {
		filter := bson.M{"_id": test_trainers1.ID}
		_, err_delete := db.Collection("trainers").DeleteOne(context.Background(), filter)
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete trainer1: %v", err_delete)
		}

		filter = bson.M{"_id": test_trainers2.ID}
		_, err_delete = db.Collection("trainers").DeleteOne(context.Background(), filter)
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete trainer2: %v", err_delete)
		}
	})

	get_trainers_request := running_trackerpb.GetTrainersRequest{
		Trainers: &running_trackerpb.Trainer{
			Brand: "test_brand1",
		},
		SortBy: []*running_trackerpb.SortField{
			{
				Field: "model",
				Order: running_trackerpb.Order_ASC,
			},
		},
	}

	// response check
	res_get, err := new_client.Trainers.GetTrainers(ctx, &get_trainers_request)
	if err != nil {
		t.Fatal("Could not get ", err)
	}
	log.Println("GET response:", res_get)

	if len(res_get.Trainers) == 0 {
		t.Fatalf("No trainers returned from GetTrainers")
	}

	expectedCount := 2
	if len(res_get.Trainers) != expectedCount {
		t.Fatalf("Expected %d trainers, got %d", expectedCount, len(res_get.Trainers))
	}

	// check sorting
	if res_get.Trainers[0].Model != test_trainers1.Model {
		t.Fatalf("Expected first trainer model to be %s, got '%s'", test_trainers1.Model, res_get.Trainers[0].Model)
	}
	if res_get.Trainers[1].Model != test_trainers2.Model {
		t.Fatalf("Expected second trainer model to be %s, got '%s'", test_trainers2.Model, res_get.Trainers[1].Model)
	}
}

func TestUpdateTrainers_ok(t *testing.T) {
	// The test calls UpdateTrainers and checks if all requested fields were updated in DB

	// db and test set up
	test_trainers := insertTestTrainer(t, mongo_client, "test_brand1", "model1")

	t.Cleanup(func() {
		_, err_delete := db.Collection("trainers").DeleteOne(context.Background(), bson.M{"_id": test_trainers.ID})
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete trainers: %v", err_delete)
		}
	})

	new_model := "Pegasus Trail 5"
	new_status := running_trackerpb.TrainerStatus_RETIRED

	update_trainers_request := running_trackerpb.UpdateTrainersRequest{
		Trainers: []*running_trackerpb.Trainer{
			{
				Id:     test_trainers.ID.Hex(),
				Model:  new_model,
				Status: new_status,
			},
		},
	}

	//response check
	res_update, err := new_client.Trainers.UpdateTrainers(ctx, &update_trainers_request)
	if err != nil {
		t.Fatalf("Could not update: %s", err)
	}
	log.Println("UPDATE response:", res_update)

	if len(res_update.Ids) == 0 {
		t.Fatalf("Update operation did not return any IDs")
	}

	//db check
	var updatedTrainer TrainerDocument
	filter := bson.M{"_id": test_trainers.ID}
	err = db.Collection("trainers").FindOne(context.Background(), filter).Decode(&updatedTrainer)
	if err != nil {
		t.Fatalf("Could not query the DB: %s ", err)
	}

	if new_model != updatedTrainer.Model {
		t.Fatalf("model not updated: expected %v, got %v", new_model, updatedTrainer.Model)
	}
	if int32(new_status) != updatedTrainer.Status {
		t.Fatalf("status not updated: expected %v, got %v", int32(new_status), updatedTrainer.Status)
	}

}

func TestDeleteTrainers_ok(t *testing.T) {
	// The test calls DeleteTrainers and validates that requested data were deleted

	test_trainers := insertTestTrainer(t, mongo_client, "test_brand1", "model1")
	t.Cleanup(func() {
		_, err_delete := db.Collection("trainers").DeleteOne(context.Background(), bson.M{"_id": test_trainers.ID})
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete trainers : %v", err_delete)
		}
	})

	delete_trainers_request := running_trackerpb.DeleteTrainersRequest{
		Ids: []string{
			test_trainers.ID.Hex(),
		},
	}

	//response check
	res_delete, err := new_client.Trainers.DeleteTrainers(ctx, &delete_trainers_request)
	if err != nil {
		t.Fatalf("Could not delete %s", err)
	}
	log.Println("DELETE response:", res_delete)

	if len(res_delete.Ids) == 0 {
		t.Fatalf("Delete operation did not return any IDs")
	}

	//db check
	filter := bson.M{"_id": test_trainers.ID}
	cursor, err := db.Collection("trainers").Find(context.Background(), filter)
	if err != nil {
		t.Fatalf("Could not query the DB: %s ", err)
	}
	var found []bson.M
	err = cursor.All(ctx, &found)
	if err != nil {
		t.Fatalf("Decoding failed: %s ", err)
	}

	if len(found) != 0 {
		t.Fatal("Delete failed, number of records found in db: ", len(found))
	}
}

func TestAddActivities_ok(t *testing.T) {
	// The test sends AddActivity request, checks the response and verifies that data were inserted in DB

	add_activities_request := running_trackerpb.AddActivitiesRequest{
		Activities: []*running_trackerpb.Activity{
			{
				Name:          "running",
				Duration:      45,
				Distance:      8,
				Date:          "2026-01-01",
				TrainersBrand: "test_brand",
				TrainersModel: "test_model",
			},
			{
				Name:     "cycling",
				Duration: 135,
				Distance: 50,
				Date:     "2026-01-02",
			},
		},
	}
	res_add, err := new_client.Activities.AddActivities(ctx, &add_activities_request)
	if err != nil {
		t.Fatal("Could not add ", err)
	}

	log.Println("IDs:", res_add.Ids)

	//db check
	objectIds := make([]primitive.ObjectID, len(res_add.Ids))
	for _, id := range res_add.Ids {
		objectId, _ := primitive.ObjectIDFromHex(id)
		objectIds = append(objectIds, objectId)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	// db clean up
	t.Cleanup(func() {
		_, err_delete := db.Collection("tracked_activities").DeleteMany(ctx, filter)
		if err_delete != nil {
			t.Logf("CLEAN_UP: Failed to delete activity: %v", err)
		}
	})

	cursor, err := db.Collection("tracked_activities").Find(context.Background(), filter)
	if err != nil {
		t.Fatalf("Could not query the DB: %s ", err)
	}
	var foundIds []bson.M
	err = cursor.All(ctx, &foundIds)
	if err != nil {
		t.Fatalf("Decoding failed: %s ", err)
	}

	if len(foundIds) != len(add_activities_request.Activities) {
		t.Fatalf("Number of found IDs in DB does not match number of added trainers")
	}

	//response check
	expected_message := "Activities were added to the database"
	if expected_message != res_add.Message {
		t.Fatalf("Expected message %s, got %s", expected_message, res_add.Message)
	}

	if len(res_add.Ids) != len(add_activities_request.Activities) {
		t.Fatalf("Number of returned IDs does not match number of added activities")
	}
}
