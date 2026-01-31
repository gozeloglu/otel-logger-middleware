/*
Package otel_logger_middleware provides a slog.Handler middleware that adds OpenTelemetry
trace information (Trace ID, Span ID, and Trace Flags) to log records.

It extracts trace context from the `context.Context` passed to the log methods and appends
the trace information as attributes to the log record. The attribute keys can be customized
using different naming conventions via `namingType`.

# Usage

To use the middleware, wrap your existing `slog.Handler` with `NewOtelLoggerMiddleware`
providing the desired naming convention.

	package main

	import (
		"context"
		"log/slog"
		"os"

		otel_logger_middleware "github.com/gozeloglu/otel-logger-middleware"
		"go.opentelemetry.io/otel"
	)

	func main() {
		// Initialize the base handler
		baseHandler := slog.NewJSONHandler(os.Stdout, nil)

		// Create the middleware with Semantic Conventions (trace.id, span.id)
		handler := otel_logger_middleware.NewOtelLoggerMiddleware(baseHandler, otel_logger_middleware.SemConv)
		logger := slog.New(handler)

		// Get a context with a span (assuming OpenTelemetry is configured)
		ctx := context.Background()
		tracer := otel.Tracer("example-tracer")
		ctx, span := tracer.Start(ctx, "example-operation")
		defer span.End()

		// Log with context to include trace info
		logger.InfoContext(ctx, "handling request", slog.String("key", "value"))
	}

# Naming Conventions

The package supports several naming conventions for trace attributes:

  - SemConv: uses dot notation (e.g., "trace.id", "span.id"). This complies with OpenTelemetry semantic conventions.
  - SnakeCase: uses underscores (e.g., "trace_id", "span_id").
  - CamelCase: uses camelCase (e.g., "traceId", "spanId").
  - PascalCase: uses PascalCase (e.g., "TraceId", "SpanId").
*/
package otel_logger_middleware
