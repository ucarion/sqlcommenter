// Package otelsqlcommenter provides a convenience function for serializing
// OpenTelemetry span context into SQLCommenter attributes.
package otelsqlcommenter

import (
	"context"
	"fmt"

	"github.com/ucarion/sqlcommenter"
	"go.opentelemetry.io/otel/trace"
)

// NewContext returns a new Context carrying "traceparent" and "tracestate"
// attributes.
//
// "traceparent" and "tracestate" will be formatted in accordance with W3C Trace
// Context encoding (https://www.w3.org/TR/trace-context/).
//
// NewContext returns ctx as-is if it does not carry a valid trace.SpanContext.
func NewContext(ctx context.Context) context.Context {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return ctx
	}

	// this code is adapted from:
	//
	// https://github.com/open-telemetry/opentelemetry-go/blob/a50cf6aadd582f9760c578e2c4b5230b6c30913d/propagation/trace_context.go#L47-L66
	return sqlcommenter.NewContext(ctx,
		"traceparent", fmt.Sprintf("%.2x-%s-%s-%s", 0, sc.TraceID(), sc.SpanID(), sc.TraceFlags()&trace.FlagsSampled),
		"tracestate", sc.TraceState().String(),
	)
}
