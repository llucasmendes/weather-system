package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

type CEPRequest struct {
	CEP string `json:"cep" binding:"required,len=8"`
}

func HandleCEP(c *gin.Context) {
	var req CEPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	ctx := c.Request.Context()
	tracer := otel.Tracer("service_a")
	_, span := tracer.Start(ctx, "HandleCEP")
	defer span.End()

	response, err := ForwardToServiceB(ctx, req.CEP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}
