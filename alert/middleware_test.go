package alert_test

import (
	"easy-api-prom-alert-sms/alert"
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/utils/file"

	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

const configContent string = `---
easy_api_prom_sms_alert:
  simulation: true
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
  recipients:
  - name: "admin"
    members:
    - "xxxxxxxxxx"
`

const body string = `
{
	"version": "4",
	"groupKey": "{alertname=\"InstanceDown\"}",
	"status": "firing",
	"receiver": "webhook",
	"groupLabels": {
	  "alertname": "InstanceDown"
	},
	"commonLabels": {
	  "alertname": "InstanceDown",
	  "job": "myjob",
	  "severity": "critical"
	},
	"commonAnnotations": {
	  "summary": "Instance xxx down",
	  "description": "The instance xxx is down."
	},
	"externalURL": "http://prometheus.example.com",
	"alerts": [
	  {
		"status": "firing",
		"labels": {
		  "alertname": "InstanceDown",
		  "instance": "localhost:9090",
		  "team": "urgence",
		  "job": "myjob",
		  "severity": "critical"
		},
		"annotations": {
		  "summary": "Instance localhost:9090 down",
		  "description": "The instance localhost:9090 is down."
		},
		"startsAt": "2023-05-20T14:28:00.000Z",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://prometheus.example.com/graph?g0.expr=up%3D%3D0&g0.tab=1"
	  }
	]
}
`

func triggerTest(t *testing.T, statusCode int, credential string, reqBody io.Reader) {
	filename, err := file.CreateConfigFileForTesting(configContent)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(filename)

	configLoaded, err := config.LoadConfig(filename, validate)
	if err != nil {
		t.Fatal(err.Error())
	}

	alertSender := alert.NewAlertSender(configLoaded)

	req, err := http.NewRequest("POST", "/api-alert", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	if credential != "" {
		req.Header.Add("Authorization", credential)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api-alert", alertSender.AlertHandler).Methods("POST")
	router.Use(alertSender.AuthMiddleware)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != statusCode {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, statusCode)
	}
}

func TestAuthentication(t *testing.T) {
	t.Run("No authorization header", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "", nil)
	})

	t.Run("Invalid authorization header", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "xxxxx", nil)
	})

	t.Run("Failed to decode base64 token", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "Basic xxxxx", nil)
	})

	t.Run("Invalid username or password", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "Basic eHh4eHg6eHh4", nil)
	})

	t.Run("Correct username and password", func(t *testing.T) {
		bodyReader := strings.NewReader(body)
		triggerTest(t, http.StatusNoContent, "Basic eHh4eHg6eHh4eHh4eHg=", bodyReader)
	})
}
