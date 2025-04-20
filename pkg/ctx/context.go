package ctx

import (
	"context"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/logr"
)

type LogContextKeyType string

const LogContextKey LogContextKeyType = "logger"

func WithLogger(ctx context.Context, logger logr.Logger) context.Context {
	return context.WithValue(ctx, LogContextKey, logger)
}

func GetLogger(ctx context.Context) logr.Logger {
	if logger, ok := ctx.Value(LogContextKey).(logr.Logger); ok {
		return logger
	}

	return logr.Logger{}
}
