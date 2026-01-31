package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	otelLoggerMiddleware "github.com/gozeloglu/otel-logger-middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.38.0"
)

var (
	tracer = otel.Tracer("Example")
)

type Runner struct {
	logger *slog.Logger
}

func main() {
	ctx := context.Background()
	shutdown, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)
	spanCtx, span := tracer.Start(ctx, "main")
	defer span.End()

	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(otelLoggerMiddleware.NewOtelLoggerMiddleware(baseHandler, otelLoggerMiddleware.SemConv))
	logger.InfoContext(spanCtx, "Logger is done")
	r := &Runner{
		logger: logger,
	}
	r.run(spanCtx)
}

func (r *Runner) run(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "run")
	defer span.End()

	r.logger.InfoContext(ctx, "run function")
	r.run2(ctx)

}

func (r *Runner) run2(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "run2")
	defer span.End()

	r.logger.ErrorContext(ctx, "run2 function")

}

func initTracer() (func(context.Context) error, error) {
	// 1. Exporter (stdout)
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	// 2. Resource (service metadata)
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("otel-logger-middleware-example"),
			attribute.String("environment", "local"),
		),
	)
	if err != nil {
		return nil, err
	}

	// 3. TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	// 4. Global provider set et
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
