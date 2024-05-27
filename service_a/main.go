package main

import (
	"context"
	"log"
	"os"
	"service_a/handler"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func main() {
	ctx := context.Background()
	tracerProvider, err := setupTracer(ctx)
	if err != nil {
		log.Fatalf("failed to setup tracer: %v", err)
	}
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer provider: %v", err)
		}
	}()

	r := gin.Default()
	r.Use(otelgin.Middleware("service_a"))

	r.POST("/cep", handler.HandleCEP)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}

func setupTracer(_ context.Context) (*trace.TracerProvider, error) {
	endpoint := "http://localhost:9411/api/v2/spans"
	exporter, err := zipkin.New(endpoint)
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service_a"),
		)),
	)
	return tp, nil
}
