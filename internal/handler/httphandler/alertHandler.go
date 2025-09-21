package httphandler

import (
	"easy-api-prom-alert-sms/internal/domain"
	"easy-api-prom-alert-sms/pkg/logger"

	"encoding/json"
	"net/http"

	"github.com/prometheus/alertmanager/template"
)

func (h *HTTPHandler) HandleHealthCheck(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandler) HandleAlert(resp http.ResponseWriter, req *http.Request) {
	var alertData template.Data

	if err := json.NewDecoder(req.Body).Decode(&alertData); err != nil {
		logger.Error("failed to parse content : %s", err.Error())
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		if err := h.Usecases.IAlertUsecase.Send(domain.Alert{Data: &alertData}); err != nil {
			logger.Error("failed to send alert : %s", err.Error())
		}
	}()

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}
