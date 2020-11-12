package providers

import (
	"os"

	"github.com/go-kit/kit/log"
)

func NewKitLoggerSet() log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger
}
