package httphandler

import (
	"easy-api-prom-alert-sms/internal/usecase"
)

type HTTPHandler struct {
	Usecases *usecase.Usecases
}

func NewHTTPHandler(usecases *usecase.Usecases) *HTTPHandler {
	return &HTTPHandler{usecases}
}
