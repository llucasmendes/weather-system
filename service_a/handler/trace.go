package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
)

type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func ForwardToServiceB(ctx context.Context, cep string) (*TemperatureResponse, error) {
	client := resty.New()
	tracer := otel.Tracer("service_a")

	var response *TemperatureResponse
	_, span := tracer.Start(ctx, "ForwardToServiceB")
	defer span.End()

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"cep": cep}).
		SetResult(&response).
		Post("http://service_b:8082/weather")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to get response from service B")
	}

	return response, nil
}
