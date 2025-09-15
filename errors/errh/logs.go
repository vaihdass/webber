package errh

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
)

const (
	errCodeKey = "xerr_error_code"
	errTypeKey = "xerr_error_type"
	errDescKey = "xerr_error_desc"
)

// A LoggingLevel is a logging priority. Higher levels are more important.
type LoggingLevel int8

const (
	// UnknownLogging no logging level.
	UnknownLogging LoggingLevel = iota
	// DebugLogging logs are typically voluminous, and are usually disabled in
	// production.
	DebugLogging
	// InfoLogging is the default logging priority.
	InfoLogging
	// WarnLogging logs are more important than Info, but don't need individual
	// human review.
	WarnLogging
	// ErrorLogging logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLogging
)

func (l LoggingLevel) toSlog() slog.Level {
	var lvl slog.Level

	switch l {
	case UnknownLogging:
		lvl = slog.LevelDebug
	case DebugLogging:
		lvl = slog.LevelDebug
	case InfoLogging:
		lvl = slog.LevelInfo
	case WarnLogging:
		lvl = slog.LevelWarn
	case ErrorLogging:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelDebug
	}

	return lvl
}

func log(
	ctx context.Context, logger *slog.Logger, lvl LoggingLevel,
	code codes.Code, errMsg, desc, errType string, kvs []any, span spanLogger,
) {
	if lvl == UnknownLogging && span == nil {
		return
	}

	kvs = append(kvs,
		errCodeKey, code.String(),
		errTypeKey, errType,
		errDescKey, desc)

	if span != nil {
		span.LogKV(kvs...)
	}

	if lvl == UnknownLogging {
		return
	}

	if logger != nil {
		logger.Log(ctx, lvl.toSlog(), errMsg, kvs...)
	}
}
