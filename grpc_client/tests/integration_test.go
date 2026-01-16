package tests

import (
	"context"
	"grpcclient/client"
	"log"
	"os"
	"testing"
	"time"

	running_trackerpb "grpcclient/proto/generated_files"

	"github.com/joho/godotenv"
)

var (
	new_client *client.Client
	ctx        context.Context
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	certPath := os.Getenv("CERT_PATH")
	token := os.Getenv("OAUTH_TOKEN")

	config := client.Config{
		CertPath:   certPath,
		OAuthToken: token,
	}

	new_client, err = client.CreateTrainersServiceClient(config)
	if err != nil {
		log.Fatal("Connection to database failed ", err)
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)

	code := m.Run()

	new_client.Conn.Close()
	cancel()
	os.Exit(code)
}

func TestAddTrainers_ok(t *testing.T) {

	add_trainers_request := running_trackerpb.AddTrainersRequest{
		Trainers: []*running_trackerpb.Trainer{
			{
				Brand:            "Nike",
				Model:            "Pegasus Trail 3",
				PurchaseDate:     "2024-02-04",
				ExpectedLifespan: 700,
				SurfaceType:      running_trackerpb.SurfaceType_ROAD_TO_TRAIL,
				Status:           running_trackerpb.TrainerStatus_NEW,
			},
			{
				Brand:            "Nike",
				Model:            "Pegasus Trail 4",
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
	log.Println("Response message:", res_add.Message)
	if len(res_add.Ids) != len(add_trainers_request.Trainers) {
		t.Fatalf("Number of returned IDs does not match number of added trainers")
	}
}

func TestGetTrainers_ok(t *testing.T) {

	get_trainers_request := running_trackerpb.GetTrainersRequest{
		Trainers: &running_trackerpb.Trainer{
			Brand: "Nike",
		},
		SortBy: []*running_trackerpb.SortField{
			{
				Field: "purchase_date",
				Order: running_trackerpb.Order_ASC,
			},
		},
	}
	res_get, err := new_client.Trainers.GetTrainers(ctx, &get_trainers_request)
	if err != nil {
		t.Fatal("Could not get ", err)
	}
	log.Println("GET response:", res_get)

	if len(res_get.Trainers) == 0 {
		t.Fatalf("No trainers returned from GetTrainers")
	}
}

func TestUpdateTrainers_ok(t *testing.T) {
	update_trainers_request := running_trackerpb.UpdateTrainersRequest{
		Trainers: []*running_trackerpb.Trainer{
			{
				Id:           "68ee6e8bdfd4eb56a49e3549",
				Brand:        "Nike",
				Model:        "Pegasus Trail 3",
				PurchaseDate: "2024-02-04",
				Status:       running_trackerpb.TrainerStatus_RETIRED,
			},
		},
	}

	res_update, err := new_client.Trainers.UpdateTrainers(ctx, &update_trainers_request)
	if err != nil {
		t.Fatalf("Could not update: %s", err)
	}
	log.Println("UPDATE response:", res_update)

	if len(res_update.Ids) == 0 {
		t.Fatalf("Update operation did not return any IDs")
	}
}

func TestDeleteTrainers_ok(t *testing.T) {
	add_trainers_request := running_trackerpb.AddTrainersRequest{
		Trainers: []*running_trackerpb.Trainer{
			{
				Brand:            "Nike",
				Model:            "Pegasus Trail 5",
				PurchaseDate:     "2025-02-04",
				ExpectedLifespan: 700,
				SurfaceType:      running_trackerpb.SurfaceType_ROAD_TO_TRAIL,
				Status:           running_trackerpb.TrainerStatus_NEW,
			},
		},
	}
	res_add, err := new_client.Trainers.AddTrainers(ctx, &add_trainers_request)
	if err != nil {
		t.Fatalf("Could not add %s", err)
	}

	delete_trainers_request := running_trackerpb.DeleteTrainersRequest{
		Ids: []string{
			res_add.Ids[0],
		},
	}
	res_delete, err := new_client.Trainers.DeleteTrainers(ctx, &delete_trainers_request)
	if err != nil {
		t.Fatalf("Could not delete %s", err)
	}
	log.Println("DELETE response:", res_delete)

	if len(res_delete.Ids) == 0 {
		t.Fatalf("Delete operation did not return any IDs")
	}

}
