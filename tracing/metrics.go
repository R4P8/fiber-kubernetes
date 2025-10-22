package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	apiMetric "go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	RequestCount    apiMetric.Int64Counter
	RequestDuration apiMetric.Float64Histogram
)

func InitMeter(ctx context.Context, serviceName, otlpEndpoint string) func(context.Context) error {

	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(otlpEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP metric exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource for metrics: %v", err)
	}

	meterProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(sdkMetric.NewPeriodicReader(exp)),
		sdkMetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	meter := otel.Meter(serviceName)

	RequestCount, err = meter.Int64Counter(
		"http_requests_total",
		apiMetric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	RequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		apiMetric.WithDescription("Duration of HTTP requests in seconds"),
	)
	if err != nil {
		log.Fatalf("failed to create histogram: %v", err)
	}

	return meterProvider.Shutdown
}
