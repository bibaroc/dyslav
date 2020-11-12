//+build wireinject

package deps

import (
	"github.com/bibaroc/dyslav/deps/providers"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
)

var (
	loggerSet = wire.NewSet(
		providers.NewKitLoggerSet,
	)
)

func InjectLogger() log.Logger {
	wire.Build(
		loggerSet,
	)

	return log.Logger(nil)
}
