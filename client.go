package login

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	status "google.golang.org/grpc/status"
)

type LoginClient struct {
	stub   LoginServiceClient
	entity string
}

// CreateLoginClient creates a new login client
func CreateLoginClient(host, entity, pathToCert string, port int) *LoginClient {
	var sdk LoginClient
	creds, err := credentials.NewClientTLSFromFile(pathToCert, host)
	if err != nil {
		log.Fatal(err)
	}
	// connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	sdk.stub = NewLoginServiceClient(conn)
	sdk.entity = entity
	return &sdk
}

// Register a new user
func (s *LoginClient) Register(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := s.stub.Register(ctx, &RegisterRequest{Entity: s.entity, Email: email, Password: password})
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.NotFound {
			return "", errors.New("entity no found")
		}
	}
	if resp != nil {
		return resp.GetToken(), err
	}

	return "", err
}
