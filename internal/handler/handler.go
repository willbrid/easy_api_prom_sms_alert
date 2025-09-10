package handler

import (
	"easy-api-prom-alert-sms/internal/entity"
	"easy-api-prom-alert-sms/internal/usecase"
	"easy-api-prom-alert-sms/pkg/logger"

	"encoding/json"
	"net/http"

	"github.com/prometheus/alertmanager/template"
)

type Handler struct {
	iAlert  usecase.IAlert
	iLogger logger.ILogger
}

func NewHandler(a usecase.IAlert, l logger.ILogger) *Handler {
	return &Handler{a, l}
}

func (c *Handler) HandleHealthCheck(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}

func (c *Handler) HandleAlert(resp http.ResponseWriter, req *http.Request) {
	var alertData template.Data

	if err := json.NewDecoder(req.Body).Decode(&alertData); err != nil {
		c.iLogger.Error("failed to parse content : %s", err.Error())
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		if err := c.iAlert.Send(entity.Alert{Data: &alertData}, c.iLogger); err != nil {
			c.iLogger.Error("failed to send alert : %s", err.Error())
		}
	}()

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}
