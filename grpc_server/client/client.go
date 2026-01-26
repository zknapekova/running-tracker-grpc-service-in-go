package client

import (
	"errors"
	running_trackerpb "grpcserver/proto/generated_files"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

type Config struct {
	CertPath   string
	OAuthToken string
}

type Client struct {
	Conn        *grpc.ClientConn
	Trainers    running_trackerpb.TrainersServiceClient
	Activities  running_trackerpb.ActivitiesServiceClient
	HealthCheck running_trackerpb.HealthCheckServiceClient
}

func CreateServiceClient(cfg Config) (*Client, error) {

	if cfg.CertPath == "" {
		return nil, errors.New("CERT_PATH not set")
	}
	token := &oauth2.Token{
		AccessToken: cfg.OAuthToken,
	}

	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(token)}
	creds, err := credentials.NewClientTLSFromFile(cfg.CertPath, "")
	if err != nil {
		log.Println("Failed to load certificate:", err)
		return nil, err
	}

	// use TLS for secure connection
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Println("Did not connect:", err)
		return nil, err
	}

	state := conn.GetState()
	log.Println("Connection State: ", state)

	//create new clients
	trainersClient := running_trackerpb.NewTrainersServiceClient(conn)
	activitiesClient := running_trackerpb.NewActivitiesServiceClient(conn)
	healthCheckClient := running_trackerpb.NewHealthCheckServiceClient(conn)

	return &Client{
		Conn:        conn,
		Trainers:    trainersClient,
		Activities:  activitiesClient,
		HealthCheck: healthCheckClient,
	}, nil
}
