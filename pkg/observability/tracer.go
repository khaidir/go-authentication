package observability

import (
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracerLocal() {
	file, err := os.Create("otel-local-traces.json")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	// exporter, err := stdouttrace.New(stdouttrace.WithWriter(os.Stdout), stdouttrace.WithPrettyPrint())
	exporter, _ := stdouttrace.New(stdouttrace.WithWriter(file))
	if err != nil {
		log.Fatalf("failed to initialize stdout exporter: %v", err)
	}
	bsp := trace.NewSimpleSpanProcessor(exporter)
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(0.5)), // 5% sampling
		trace.WithSpanProcessor(bsp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("auth-service"),
		)),
	)

	otel.SetTracerProvider(tp)
}
