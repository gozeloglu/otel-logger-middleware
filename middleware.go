package otel_logger_middleware

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type namingConverter interface {
	convert(span trace.Span, record *slog.Record)
}

type otelLoggerMiddleware struct {
	slog.Handler
	converter namingConverter
}

type namingType int

const (
	// SemConv stands for OpenTelemetry semantic convention. Example, "trace.id"
	SemConv namingType = iota

	// SnakeCase adds underscore. Example, "trace_id"
	SnakeCase

	// CamelCase makes second word's first letter uppercase. Example, "camelCase"
	CamelCase

	// PascalCase makes first and second word's first letters uppercase. Example, "CamelCase"
	PascalCase
)

func NewOtelLoggerMiddleware(baseHandler slog.Handler, converter namingConverter) slog.Handler {
	return &otelLoggerMiddleware{
		Handler:   baseHandler,
		converter: converter,
	}
}

func (o *otelLoggerMiddleware) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() && span.SpanContext().HasSpanID() && span.SpanContext().HasTraceID() {
		o.converter.convert(span, &record)
	}

	return o.Handler.Handle(ctx, record)
}

func (t namingType) convert(span trace.Span, record *slog.Record) {
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()
	traceFlags := span.SpanContext().TraceFlags().String()
	switch t {
	case SemConv:
		record.Add(
			"span.id", spanId,
			"trace.id", traceId,
			"trace.flags", traceFlags,
		)
		return
	case SnakeCase:
		record.Add(
			"span_id", spanId,
			"trace_id", traceId,
			"trace_flags", traceFlags,
		)
		return
	case CamelCase:
		record.Add(
			"spanId", spanId,
			"traceId", traceId,
			"traceFlags", traceFlags,
		)
		return
	case PascalCase:
		record.Add(
			"SpanId", spanId,
			"TraceId", traceId,
			"TraceFlags", traceFlags,
		)
		return
	}
}
