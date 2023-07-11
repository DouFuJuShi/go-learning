# OpenTelemetry

```go

import (
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"

    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"

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
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
        sdktrace.WithBatcher(exporter), 
        sdktrace.WithResource(resources), 
    )

    // SetTracerProvider "go.opentelemetry.io/otel"
    otel.SetTracerProvider(tp)

    // tracer
    tracerName := "tracer"
    var tracerOption []trace.TracerOption
    tracer := otel.GetTracerProvider().Tracer(tracerName, tracerOption...)

    // root span
    rootName := "root_span"
    var rootStartOption []trace.SpanStartOption
    parentCtx, rootSpan := tracer.Start(context.Background(), rootName, rootStartOption...)
    defer rootSpan.End()

    // child span
    childName := "child_span"
    var childStartOption []trace.SpanStartOption
    childCtx, childSpan := tracer.Start(parentCtx, childName, childStartOption...)

    // set Attributes to span
    childSpan.SetAttributes(attribute.String("controller", "books"))
}
```
