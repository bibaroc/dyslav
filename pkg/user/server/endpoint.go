package server

import (
	"context"
	"time"

	pb "github.com/bibaroc/dyslav/pkg/user/proto"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Set struct {
	CreateUser endpoint.Endpoint
	Liveness   endpoint.Endpoint
	Readiness  endpoint.Endpoint
}

func MakeEndpointSet(srv pb.UserServiceServer) Set {
	var createUser endpoint.Endpoint
	{
		createUser = MakeCreateUserEndpoint(srv)
		createUser = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(createUser)
		createUser = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createUser)
	}

	var liveness endpoint.Endpoint
	{
		liveness = MakeLivenessEndpoint(srv)
		liveness = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(liveness)
		liveness = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(liveness)
	}

	var readiness endpoint.Endpoint
	{
		readiness = MakeReadinessEndpoint(srv)
		readiness = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(readiness)
		readiness = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(readiness)
	}

	return Set{
		CreateUser: createUser,
		Liveness:   liveness,
		Readiness:  readiness,
	}
}

func MakeCreateUserEndpoint(srv pb.UserServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.CreateUserRequest)

		v, err := srv.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}

		return v, nil
	}
}

func MakeLivenessEndpoint(srv pb.UserServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.LivenessRequest)

		v, err := srv.Liveness(ctx, req)
		if err != nil {
			return nil, err
		}

		return v, nil
	}
}

func MakeReadinessEndpoint(srv pb.UserServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.ReadinessRequest)

		v, err := srv.Readiness(ctx, req)
		if err != nil {
			return nil, err
		}

		return v, nil
	}
}
