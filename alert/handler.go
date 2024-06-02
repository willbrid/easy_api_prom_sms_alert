package alert

import (
	"easy-api-prom-alert-sms/logging"

	"encoding/json"
	"net/http"

	"github.com/prometheus/alertmanager/template"
)

func (alertSender *AlertSender) AlertHandler(resp http.ResponseWriter, req *http.Request) {
	var alertData template.Data

	if err := json.NewDecoder(req.Body).Decode(&alertData); err != nil {
		logging.Log(logging.Error, "failed to parse content : %s", err.Error())
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	alertSender.setData(&alertData)

	go func() {
		if err := alertSender.sendAlert(); err != nil {
			logging.Log(logging.Error, "failed to send alert : %s", err.Error())
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
	}()

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}
