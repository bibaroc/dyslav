package server

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(endpoints Set, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	m.Handle("/create_user", httptransport.NewServer(
		endpoints.CreateUser,
		decodeHTTPCreateUserRequest,
		encodeHTTPGenericResponse,
		options...,
	))

	return m
}
