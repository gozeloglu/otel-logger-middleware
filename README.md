# otel-logger-middleware

`otel-logger-middleware` is a Go library that provides a middleware for `log/slog` to automatically inject OpenTelemetry trace context (Trace ID, Span ID, and Trace Flags) into your logs.

## Why is this important?

In distributed systems, correlating logs with traces is crucial for observability. When you have a log message, you often want to know which request (trace) generated it. Conversely, when looking at a trace, you want to see the associated logs.

This middleware bridges that gap by adding `trace_id`, `span_id`, and `trace_flags` to `slog.Record` automatically when a valid span is found in the context.

## Installation

```bash
go get github.com/gozeloglu/otel-logger-middleware
```

## Usage

To use this middleware, wrap your base `slog.Handler` with `otelLoggerMiddleware.NewOtelLoggerMiddleware`.

```go
package main

import (
	"context"
	"log/slog"
	"os"

	otelLoggerMiddleware "github.com/gozeloglu/otel-logger-middleware"
	"go.opentelemetry.io/otel"
)

func main() {
	// 1. Create your base handler (e.g., JSON or Text)
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// 2. Wrap it with OtelLoggerMiddleware
	// You can choose the naming convention for the keys:
	// - SemConv: "trace.id", "span.id"
	// - SnakeCase: "trace_id", "span_id"
	// - CamelCase: "traceId", "spanId"
	// - PascalCase: "TraceId", "SpanId"
	middleware := otelLoggerMiddleware.NewOtelLoggerMiddleware(baseHandler, otelLoggerMiddleware.SemConv)

	// 3. Create the logger
	logger := slog.New(middleware)

	// 4. Use the logger with a context containing a span
	ctx := context.Background()
	tracer := otel.Tracer("example-tracer")
	ctx, span := tracer.Start(ctx, "my-operation")
	defer span.End()

	// The log output will now contain trace.id and span.id
	logger.InfoContext(ctx, "This log is correlated with a trace")
}
```

## Naming Conventions

The library supports different key naming conventions to match your logging backend requirements:

- `otelLoggerMiddleware.SemConv`: Uses OpenTelemetry semantic conventions (e.g., `trace.id`, `span.id`).
- `otelLoggerMiddleware.SnakeCase`: Uses snake_case (e.g., `trace_id`, `span_id`).
- `otelLoggerMiddleware.CamelCase`: Uses camelCase (e.g., `traceId`, `spanId`).
- `otelLoggerMiddleware.PascalCase`: Uses PascalCase (e.g., `TraceId`, `SpanId`).
