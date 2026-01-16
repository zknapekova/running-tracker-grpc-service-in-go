package client

import (
	running_trackerpb "grpcclient/proto/generated_files"
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
	Conn     *grpc.ClientConn
	Trainers running_trackerpb.TrainersServiceClient
}

func CreateTrainersServiceClient(cfg Config) (*Client, error) {

	if cfg.CertPath == "" {
		log.Fatal("CERT_PATH not set")
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

	//create new client
	client := running_trackerpb.NewTrainersServiceClient(conn)

	return &Client{
		Conn:     conn,
		Trainers: client,
	}, nil
}
