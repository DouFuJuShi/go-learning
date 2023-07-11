# OpenTelemetry

```go

import (
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"

    semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var collectorURL = "collector endpoint"

func main() {
    // exporter go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
    secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
    if len(insecure) > 0 {
        secureOption = otlptracegrpc.WithInsecure()
    }

    // otlptrace go.opentelemetry.io/otel/exporters/otlp/otlptrace
    exporter, err := otlptrace.New(
        context.Background(),
        otlptracegrpc.NewClient(
            secureOption,
            otlptracegrpc.WithEndpoint(collectorURL),
        ),
    )

    if err != nil {
        log.Fatal(err)
    }


    // Resource go.opentelemetry.io/otel/sdk/resource
    resources := resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceNameKey.String("myService"),
        semconv.ServiceVersionKey.String("1.0.0"),
        semconv.ServiceInstanceIDKey.String("abcdef12345"),
    )

    // TracerProvider go.opentelemetry.io/otel/sdk/trace
    provider := sdktrace.NewTracerProvider(
        sdktrace.WithSampler(sdktrace.AlwaysSample()), // 采样设置
        sdktrace.WithBatcher(exporter), // exporter设置
        sdktrace.WithResource(resources), 
    )

    // SetTracerProvider "go.opentelemetry.io/otel"
    otel.SetTracerProvider(tp)
}
```
