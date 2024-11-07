package log

import "context"

type (
	ctxMarker struct{}

	ctxLogger struct {
		logger *Logger
	}
)

// nolint:gochecknoglobals
var ctxMarkerKey = &ctxMarker{}

// Extract takes the call-scoped Logger.
func Extract(ctx context.Context) *Logger {
	l, ok := ctx.Value(ctxMarkerKey).(*ctxLogger)
	if !ok || l == nil {
		return New(false, false, false)
	}

	return l.logger
}

// ToContext adds the log.Logger to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, logger *Logger) context.Context {
	l := &ctxLogger{
		logger: logger,
	}

	return context.WithValue(ctx, ctxMarkerKey, l)
}

// AddField adds a field to the context logger.
func AddField(ctx context.Context, key string, value any) {
	Extract(ctx).WithField(key, value)
}