package main

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	running_trackerpb "grpcclient/proto/generated_files"
	"log"
	"os"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	certPath := os.Getenv("CERT_PATH")

	if certPath == "" {
		log.Fatal("CERT_PATH not set")
	}
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		log.Fatal("Failed to load certificate:", err)
	}

	// establish insecure connection for now
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal("Did not connect:", err)
	}
	defer conn.Close()

	//create new client
	client := running_trackerpb.NewTrainersServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	//create request
	request := running_trackerpb.AddTrainersRequest{
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
	//get response
	res, err := client.AddTrainers(ctx, &request)
	if err != nil {
		log.Fatal("Could not add", err)
	}
	state := conn.GetState()
	log.Println("Connection State: ", state)

	log.Println("IDs:", res.Ids)
	log.Println("Response message:", res.Message)


}
