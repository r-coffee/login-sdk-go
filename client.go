package login

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/r-coffee/login-sdk-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	connectTimeout = 5 * time.Second
	requestTimeout = 3 * time.Second
)

type LoginClient struct {
	stub   proto.LoginServiceClient
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
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	sdk.stub = proto.NewLoginServiceClient(conn)
	sdk.entity = entity
	return &sdk
}

// Register a new user
// returns a jwt token and error
func (s *LoginClient) Register(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := s.stub.Register(ctx, &proto.RegisterRequest{Entity: s.entity, Email: email, Password: password})
	if resp != nil {
		return resp.GetToken(), err
	}

	return "", err
}

// Login a user
// returns a jwt token and error
func (s *LoginClient) Login(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := s.stub.Login(ctx, &proto.LoginRequest{Entity: s.entity, Email: email, Password: password})
	if resp != nil {
		return resp.GetToken(), err
	}

	return "", err
}

// Validate a jwt token
// returns email and error
func (s *LoginClient) Validate(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := s.stub.Validate(ctx, &proto.ValidateRequest{Entity: s.entity, Token: token})
	if resp != nil {
		return resp.GetEmail(), err
	}

	return "", err
}

// List all entities
func (s *LoginClient) List() ([]*proto.EntityTuple, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := s.stub.List(ctx, &proto.ListRequest{})
	if resp != nil {
		return resp.Entities, err
	}

	return []*proto.EntityTuple{}, err
}

// Remove an entity
func (s *LoginClient) Remove(guid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := s.stub.Remove(ctx, &proto.RemoveRequest{Guid: guid})
	return err
}
