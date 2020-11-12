package server

import (
	"context"
	"fmt"

	pb "github.com/bibaroc/dyslav/pkg/user/proto"
)

type srv struct {
}

func (s srv) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
	fmt.Println(req)
	return &pb.CreateUserResponse{}, nil
}

func NewUserServer() pb.UserServiceServer {
	return srv{}
}
