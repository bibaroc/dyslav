package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bibaroc/dyslav/deps"
	"github.com/bibaroc/dyslav/pkg"
	"github.com/bibaroc/dyslav/pkg/user/server"
	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	logger := deps.InjectLogger()

	_ = logger.Log("version", pkg.Version, "gitCommit", pkg.Commit)

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
		endpoints   = server.MakeEndpointSet(usersrv)
		httpHandler = server.NewHTTPHandler(endpoints, logger)

		g group.Group
	)

	var (
		metricsAddr = "0.0.0.0:8000"
		httpAddr    = "0.0.0.0:8001"
	)

	{ // metrics listen
		switch envStr("METRICS_PROVIDER", "") {
		case "PROMETHEUS":
			debugListener, err := net.Listen("tcp", metricsAddr)
			if err != nil {
				_ = logger.Log("transport", "metrics/HTTP", "during", "Listen", "err", err)
				os.Exit(1)
			}
			g.Add(func() error {
				_ = logger.Log("transport", "metrics/HTTP", "addr", metricsAddr)
				return http.Serve(debugListener, http.DefaultServeMux)
			}, func(error) {
				debugListener.Close()
			})
		}
	}
	{ // http user service
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			_ = logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{ // graceful shutdown
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	_ = logger.Log("exit", g.Run())
}
