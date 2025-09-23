package detachctx

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// NewDetachedContext creates a new context that is detached from the passed one but has access to its values.
//
// First, the value is searched for in new (child) context, then, if not found, in old (parent) context.
// Passed context timeout is not used (cancel too).
//
// This constructor is useful when you need to run a background process that requires information from the original context
// (for example, query metadata).
func NewDetachedContext(ctx context.Context) context.Context {
	return newComposedContext(ctx, opentracing.ContextWithSpan(context.Background(), opentracing.SpanFromContext(ctx)))
}

func newComposedContext(ctxWithValues, newCtx context.Context) context.Context {
	return &composedContext{Context: newCtx, ctxValues: ctxWithValues}
}

type composedContext struct {
	context.Context

	ctxValues context.Context
}

func (c *composedContext) Value(key any) any {
	if v := c.Context.Value(key); v != nil {
		return v
	}

	return c.ctxValues.Value(key)
}
