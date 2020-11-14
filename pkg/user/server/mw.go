package server

import (
	"context"
	"fmt"
	"time"

	pb "github.com/bibaroc/dyslav/pkg/user/proto"
	"github.com/go-kit/kit/log"
)

func StdOutLogging(
	logger log.Logger,

) func(pb.UserServiceServer) pb.UserServiceServer {
	return func(next pb.UserServiceServer) pb.UserServiceServer {
		return stdOutLoggingMW{
			logger: logger,
			next:   next,
		}
	}
}

type stdOutLoggingMW struct {
	next   pb.UserServiceServer
	logger log.Logger
}

func (s stdOutLoggingMW) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (res *pb.CreateUserResponse, err error) {
	defer func() {
		_ = s.logger.Log(
			"method", "CreateUser",
			"username", req.Username,
			"email", req.Email,
			"err", err)
	}()

	res, err = s.next.CreateUser(ctx, req)

	return
}

func (s stdOutLoggingMW) Liveness(
	ctx context.Context,
	req *pb.LivenessRequest,
) (res *pb.LivenessResponse, err error) {
	defer func() {
		_ = s.logger.Log(
			"method", "Liveness",
			"err", err)
	}()

	res, err = s.next.Liveness(ctx, req)

	return
}

func (s stdOutLoggingMW) Readiness(
	ctx context.Context,
	req *pb.ReadinessRequest,
) (res *pb.ReadinessResponse, err error) {
	defer func() {
		_ = s.logger.Log(
			"method", "Readiness",
			"err", err)
	}()

	res, err = s.next.Readiness(ctx, req)

	return
}

func PrometheusMetrics(
	metrics Metrics,
) func(pb.UserServiceServer) pb.UserServiceServer {
	return func(next pb.UserServiceServer) pb.UserServiceServer {
		return prometheusMetricsCollector{
			Metrics: metrics,
			next:    next,
		}

	}
}

type prometheusMetricsCollector struct {
	Metrics
	next pb.UserServiceServer
}

func (p prometheusMetricsCollector) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (res *pb.CreateUserResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateUser", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res, err = p.next.CreateUser(ctx, req)

	return
}

func (p prometheusMetricsCollector) Liveness(
	ctx context.Context,
	req *pb.LivenessRequest,
) (res *pb.LivenessResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Liveness", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res, err = p.next.Liveness(ctx, req)

	return
}

func (p prometheusMetricsCollector) Readiness(
	ctx context.Context,
	req *pb.ReadinessRequest,
) (res *pb.ReadinessResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Readiness", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res, err = p.next.Readiness(ctx, req)

	return
}
