package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"strings"
	"time"
)

var (
	serviceName  = "Parcels Service"
	collectorURL = "127.0.0.1:4317"
)

func initTracer() func(context.Context) error {

	//Setting the Service name from the environmental variable if exists
	if strings.TrimSpace(os.Getenv("SERVICE_NAME")) != "" {
		serviceName = os.Getenv("SERVICE_NAME")
	}
	//Setting the Collector endpoint from the environmental variable if exists
	if strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) != "" {
		collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}

	//Setting up the exporter for the tracer
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)
	//Log a fatal error if exporter could not be setup
	if err != nil {
		log.Fatal(err)
	}

	// Setting up the resources for the tracer. this includes the context and other attributes
	//to identify the source of the traces
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("language", "go"),
		),
	)
	if err != nil {
		log.Println("Could not set resources: ", err)
	}

	//Using the resources and exporter to set up a trace provider
	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

const (
	tracerKey  = "otel-go-contrib-tracer"
	tracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// ChildSpan A utility function to create child spans for specific opterations
func ChildSpan(c *gin.Context, action, name string, task func() error) error {
	//Setting up a tracer either from the existing context or creating a new one
	var tracer trace.Tracer
	tracerInterface, ok := c.Get(tracerKey)
	if ok {
		tracer, ok = tracerInterface.(trace.Tracer)
	}
	if !ok {
		tracer = otel.GetTracerProvider().Tracer(
			tracerName,
			trace.WithInstrumentationVersion(otelgin.Version()),
		)
	}
	savedContext := c.Request.Context()
	defer func() {
		c.Request = c.Request.WithContext(savedContext)
	}()
	//Adding attributes to identify the operation captured in this span
	opt := trace.WithAttributes(attribute.String("service.action", action))
	_, span := tracer.Start(savedContext, name, opt)

	// Simulate delay in operation
	time.Sleep(time.Millisecond * time.Duration(randomDelay(400, 800)))

	//running function provided in span to add to trace
	if err := task(); err != nil {
		// recording an error into the span if there is any
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("action %s failure", action))
		span.End()
		return err
	}
	//ending the span
	span.End()
	return nil
}
