package weather

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
)

type CEPRequest struct {
	CEP string `json:"cep" binding:"required,len=8"`
}

type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func GetWeather(c *gin.Context) {
	var req CEPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	ctx := c.Request.Context()
	tracer := otel.Tracer("service_b")
	_, span := tracer.Start(ctx, "GetWeather")
	defer span.End()

	city, err := getCityFromZipcode(ctx, req.CEP)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			c.JSON(http.StatusNotFound, gin.H{"message": "can not find zipcode"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	tempC, err := getTemperature(ctx, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting temperature"})
		return
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	response := TemperatureResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	c.JSON(http.StatusOK, response)
}

func getCityFromZipcode(ctx context.Context, zipcode string) (string, error) {
	client := resty.New()
	tracer := otel.Tracer("service_b")

	var city string
	_, span := tracer.Start(ctx, "getCityFromZipcode")
	defer span.End()

	resp, err := client.R().
		SetContext(ctx).
		SetResult(&map[string]interface{}{}).
		Get("https://viacep.com.br/ws/" + zipcode + "/json/")
	if err != nil {
		return "", err
	}

	data := *(resp.Result().(*map[string]interface{}))
	if _, ok := data["erro"]; ok {
		return "", errors.New("can not find zipcode")
	}

	city, exists := data["localidade"].(string)
	if !exists {
		return "", errors.New("invalid response from viaCEP")
	}

	return city, nil
}

func getTemperature(ctx context.Context, city string) (float64, error) {
	client := resty.New()
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		return 0, errors.New("missing weather API key")
	}
	tracer := otel.Tracer("service_b")

	var tempC float64
	_, span := tracer.Start(ctx, "getTemperature")
	defer span.End()

	resp, err := client.R().
		SetContext(ctx).
		SetQueryParam("key", weatherAPIKey).
		SetQueryParam("q", city).
		SetResult(&map[string]interface{}{}).
		Get("http://api.weatherapi.com/v1/current.json")
	if err != nil {
		return 0, err
	}

	data := *(resp.Result().(*map[string]interface{}))
	current, exists := data["current"].(map[string]interface{})
	if !exists {
		return 0, errors.New("invalid response from WeatherAPI")
	}

	tempC, exists = current["temp_c"].(float64)
	if !exists {
		return 0, errors.New("invalid temperature data")
	}

	return tempC, nil
}
