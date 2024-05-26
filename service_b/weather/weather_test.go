package weather

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCityFromZipcode(t *testing.T) {
	ctx := context.Background()
	city, err := getCityFromZipcode(ctx, "01001000")
	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", city)
}

func TestGetTemperature(t *testing.T) {
	ctx := context.Background()
	temp, err := getTemperature(ctx, "São Paulo")
	assert.NoError(t, err)
	assert.NotZero(t, temp)
}
