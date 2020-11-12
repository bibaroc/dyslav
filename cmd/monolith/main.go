package main

import (
	"github.com/bibaroc/dyslav/deps"
	"github.com/bibaroc/dyslav/pkg/user/server"
	"github.com/go-kit/kit/log"
)

//nolint
var (
	version string = "dev"
	commit  string = "dev"
)

func main() {
	logger := deps.InjectLogger()

	usersrv := server.NewUserServer()
	{
		switch envStr("METRICS_PROVIDER", "") {
		case "PROMETHEUS":
			usersrv = server.PrometheusMetrics(server.MakeMetrics())(usersrv)
		}

		switch envStr("LOG_PROVIDER", "") {
		case "STDOUT":
			usersrv = server.StdOutLogging(log.With(logger, "component", "usersrv"))(usersrv)
		}
	}

	var (
		endpoints = server.MakeEndpointSet(usersrv)
	)
	_ = logger.Log("version", version, "gitCommit", commit)
}
