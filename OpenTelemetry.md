# OpenTelemetry

![流程](https://github.com/DouFuJuShi/learning/assets/13195119/b19f8c4d-e752-422e-b179-10117f6c3deb)

```go
// https://signoz.io/docs/instrumentation/golang/
import (
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"

    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    // "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"

    semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

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

## Fib Demo
```go
package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
)

// name is the Tracer name used to identify this instrumentation library.
const name = "fib"

// Fibonacci returns the n-th fibonacci number.
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}

// App is a Fibonacci computation application.
type App struct {
	r io.Reader
	l *log.Logger
}

// NewApp returns a new App.
func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

// Run starts polling users for Fibonacci number requests and writes results.
func (a *App) Run(ctx context.Context) error {
	for {
		// Each execution of the run loop, we should get a new "root" span and context.
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")

		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}

		a.Write(newCtx, n)
		span.End()
	}
}

// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, err
}

// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		f, err := Fibonacci(n)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return f, err
	}(ctx)
	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

func main() {
	l := log.New(os.Stdout, "", 0)

	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
```

## Context Propagation

https://opentelemetry.io/docs/concepts/signals/baggage/    

https://uptrace.dev/opentelemetry/distributed-tracing.html#context-propagation    

[Passthrough Example](https://github.com/open-telemetry/opentelemetry-go/tree/main/example/passthrough)

## Opentelemetry-Go-Contrib
[opentelemetry-go-contrib](https://github.com/open-telemetry/opentelemetry-go-contrib)

[uptrace](https://uptrace.dev/get/instrument/)

## Docs
See ==> https://github.com/open-telemetry/docs-cn

## Distributed Tracing Tools 
See ==> https://uptrace.dev/blog/distributed-tracing-tools.html
