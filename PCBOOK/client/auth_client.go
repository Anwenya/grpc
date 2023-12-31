package client

import (
	"context"
	"time"

	"com.wlq/pcbook/pb"
	"google.golang.org/grpc"
)

type AuthClient struct {
	sevice   pb.AuthServiceClient
	username string
	password string
}

func NewAuthClient(cc *grpc.ClientConn, username string, password string) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{
		sevice:   service,
		username: username,
		password: password,
	}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.sevice.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
