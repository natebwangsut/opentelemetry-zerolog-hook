package otelzerolog

import (
	"time"

	otel "github.com/agoda-com/opentelemetry-logs-go/logs" // use otel so that when otel is stable, we can just change the import path
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "otelzerolog"
	version             = "0.0.1"
)

var instrumentationScope = instrumentation.Scope{
	Name:      instrumentationName,
	Version:   version,
	SchemaURL: semconv.SchemaURL,
}

type Hook struct {
	otel.Logger
}

var _ zerolog.Hook = (*Hook)(nil)

func NewHook(loggerProvider otel.LoggerProvider) *Hook {
	logger := loggerProvider.Logger(
		instrumentationScope.Name,
		otel.WithInstrumentationVersion(instrumentationScope.Version),
	)
	return &Hook{logger}
}

func (h Hook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if !e.Enabled() {
		return
	}

	ctx := e.GetCtx()
	span := trace.SpanFromContext(ctx).SpanContext()

	var spanID trace.SpanID
	var traceID trace.TraceID
	var traceFlags trace.TraceFlags
	if span.IsValid() {
		spanID = span.SpanID()
		traceID = span.TraceID()
		traceFlags = span.TraceFlags()
	}

	severityText := level.String()
	severityNumber := otelLevel(level)

	now := time.Now()
	lrc := otel.LogRecordConfig{
		Timestamp:            &now,
		ObservedTimestamp:    now,
		TraceId:              &traceID,
		SpanId:               &spanID,
		TraceFlags:           &traceFlags,
		SeverityText:         &severityText,
		SeverityNumber:       &severityNumber,
		Body:                 &msg,
		Resource:             nil,
		InstrumentationScope: &instrumentationScope,
		// Attributes:           &attributes, // TODO: add attributes
	}

	r := otel.NewLogRecord(lrc)
	h.Emit(r)
}

func otelLevel(level zerolog.Level) otel.SeverityNumber {
	switch level {
	case zerolog.TraceLevel:
		return otel.TRACE
	case zerolog.DebugLevel:
		return otel.DEBUG
	case zerolog.InfoLevel:
		return otel.INFO
	case zerolog.WarnLevel:
		return otel.WARN
	case zerolog.ErrorLevel:
		return otel.ERROR
	case zerolog.PanicLevel:
		return otel.ERROR
	case zerolog.FatalLevel:
		return otel.FATAL
	}
	return otel.TRACE
}
