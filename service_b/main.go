package main

import (
	"context"
	"log"
	"os"
	"service_b/weather"

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
	r.Use(otelgin.Middleware("service_b"))

	r.POST("/weather", weather.GetWeather)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}

func setupTracer(_ context.Context) (*trace.TracerProvider, error) {
	endpoint := "http://zipkin:9411/api/v2/spans" // alterado para usar o nome do serviço no Docker Compose
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
