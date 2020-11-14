package server

import (
	"context"
	"fmt"

	"github.com/bibaroc/dyslav/pkg"
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

func (s srv) Liveness(
	_ context.Context,
	_ *pb.LivenessRequest,
) (*pb.LivenessResponse, error) {
	return &pb.LivenessResponse{
		Commit:  pkg.Commit,
		Version: pkg.Version,
		Status:  pb.Status_UP,
	}, nil
}

func (s srv) Readiness(
	_ context.Context,
	_ *pb.ReadinessRequest,
) (*pb.ReadinessResponse, error) {
	return &pb.ReadinessResponse{
		Status: pb.Status_UP,
	}, nil
}

func NewUserServer() pb.UserServiceServer {
	return srv{}
}
