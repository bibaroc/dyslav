package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	pb "github.com/bibaroc/dyslav/pkg/user/proto"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(endpoints Set, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	{
		m.Handle("/create_user", httptransport.NewServer(
			endpoints.CreateUser,
			decodeHTTPCreateUserRequest,
			encodeHTTPGenericResponse,
			options...,
		))
		m.Handle("/liveness", httptransport.NewServer(
			endpoints.Liveness,
			decodeHTTPLivenessRequest,
			encodeHTTPGenericResponse,
			options...,
		))
		m.Handle("/readiness", httptransport.NewServer(
			endpoints.Readiness,
			decodeHTTPReadinessRequest,
			encodeHTTPGenericResponse,
			options...,
		))
	}

	return m
}

func decodeHTTPCreateUserRequest(
	_ context.Context,
	r *http.Request) (interface{}, error) {
	var req pb.CreateUserRequest
	err := jsonDecodeValidate(r.Body, &req)
	return &req, err
}
func decodeHTTPLivenessRequest(
	_ context.Context,
	r *http.Request) (interface{}, error) {
	return &pb.LivenessRequest{}, nil
}
func decodeHTTPReadinessRequest(
	_ context.Context,
	r *http.Request) (interface{}, error) {
	return &pb.ReadinessRequest{}, nil
}

func encodeHTTPGenericResponse(
	_ context.Context,
	w http.ResponseWriter,
	response interface{},
) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func jsonDecodeValidate(db io.ReadCloser, v interface{}) error {
	err := json.NewDecoder(db).Decode(v)
	if err != nil {
		return err
	}
	type Validator interface {
		Validate() error
	}
	validable, ok := v.(Validator)
	if ok {
		err = validable.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
