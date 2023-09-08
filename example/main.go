package main

import (
	"context"
	"os"

	otel "github.com/agoda-com/opentelemetry-logs-go"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	otelzerolog "github.com/natebwangsut/opentelemetry-zerolog-hook"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func newResource() *resource.Resource {
	host, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otlplogs-example"),
		semconv.ServiceVersion("0.0.1"),
		semconv.HostName(host),
	)
}

func main() {
	ctx := context.Background()

	exporter, _ := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogshttp.NewClient()))
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(newResource()),
	)
	otel.SetLoggerProvider(loggerProvider)

	hook := otelzerolog.NewHook(loggerProvider)
	log := log.Hook(hook)

	log.Info().Ctx(ctx).Msg("Hello OpenTelemetry")
}
