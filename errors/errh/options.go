package errh

// spanLogger represents a span (active, un-finished) in the OpenTracing system for logging purpose.
type spanLogger interface {
	LogKV(kvs ...any)
}

// handleOpts contains options for error handling.
type handleOpts struct {
	fallbackMsg string
	values      []any
	span        spanLogger
}

// Option is a function that configures handleOpts.
type Option func(*handleOpts)

// configureOptions applies the given options to handleOpts.
func configureOptions(opts ...Option) handleOpts {
	var options handleOpts

	if len(opts) == 0 {
		return options
	}

	for i := range opts {
		if opts[i] == nil {
			continue
		}

		opts[i](&options)
	}

	return options
}

// Msg sets the fallback message for unknown errors (not xerr.Error).
// Used only if handler's NotXerrCallback applied to unknown error with false result.
func Msg(fallbackMsg string) Option {
	return func(o *handleOpts) {
		o.fallbackMsg = fallbackMsg
	}
}

// Values sets key-value pairs for logging in logging and tracing.
func Values(args ...any) Option {
	return func(o *handleOpts) {
		if len(args)%2 != 0 || len(args) == 0 {
			// invalid number of arguments
			return
		}

		o.values = append(o.values, args...)
	}
}

// Span sets a span for logging values to the traces.
func Span(span spanLogger) Option {
	return func(o *handleOpts) {
		o.span = span
	}
}
