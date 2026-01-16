package main

import (
	"context"
	running_trackerpb "grpcclient/proto/generated_files"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func assert(condition bool, msg string) {
	if !condition {
		log.Fatal("TEST FAILED: ", msg)
	}
}

func add_trainers_ok(ctx context.Context, client running_trackerpb.TrainersServiceClient) {

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
	res_add, err := client.AddTrainers(ctx, &add_trainers_request)
	if err != nil {
		log.Fatal("Could not add", err)
	}

	log.Println("IDs:", res_add.Ids)
	log.Println("Response message:", res_add.Message)
	assert(len(res_add.Ids) == len(add_trainers_request.Trainers), "Number of returned IDs does not match number of added trainers")
}

func get_trainers_ok(ctx context.Context, client running_trackerpb.TrainersServiceClient) {

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
	res_get, err := client.GetTrainers(ctx, &get_trainers_request)
	if err != nil {
		log.Fatal("Could not get", err)
	}
	log.Println("GET response:", res_get)
	assert(len(res_get.Trainers) > 0, "No trainers returned from GetTrainers")
}

func update_trainers_ok(ctx context.Context, client running_trackerpb.TrainersServiceClient) {
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

	res_update, err := client.UpdateTrainers(ctx, &update_trainers_request)
	if err != nil {
		log.Fatal("Could not update", err)
	}
	log.Println("UPDATE response:", res_update)
	assert(len(res_update.Ids) > 0, "Update operation did not return any IDs")
}

func delete_trainers_ok(ctx context.Context, client running_trackerpb.TrainersServiceClient) {
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
	res_add, err := client.AddTrainers(ctx, &add_trainers_request)
	if err != nil {
		log.Fatal("Could not add", err)
	}

	delete_trainers_request := running_trackerpb.DeleteTrainersRequest{
		Ids: []string{
			res_add.Ids[0],
		},
	}
	res_delete, err := client.DeleteTrainers(ctx, &delete_trainers_request)
	if err != nil {
		log.Fatal("Could not delete", err)
	}
	log.Println("DELETE response:", res_delete)
	assert(len(res_delete.Ids) > 0, "Delete operation did not return any IDs")

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	certPath := os.Getenv("CERT_PATH")
	if certPath == "" {
		log.Fatal("CERT_PATH not set")
	}
	token := &oauth2.Token{
		AccessToken: os.Getenv("OAUTH_TOKEN"),
	}

	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(token)}
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		log.Fatal("Failed to load certificate:", err)
	}

	// use TLS for secure connection
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatal("Did not connect:", err)
	}
	defer conn.Close()

	state := conn.GetState()
	log.Println("Connection State: ", state)

	//create new client
	client := running_trackerpb.NewTrainersServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	add_trainers_ok(ctx, client)
	get_trainers_ok(ctx, client)
	update_trainers_ok(ctx, client)
	delete_trainers_ok(ctx, client)
}
