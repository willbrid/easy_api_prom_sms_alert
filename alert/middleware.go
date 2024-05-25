package alert

import (
	"easy-api-prom-alert-sms/logging"

	"encoding/base64"
	"net/http"
	"strings"
)

func (alertSender *AlertSender) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var auth string = req.Header.Get("Authorization")

		if alertSender.config.EasyAPIPromAlertSMS.Auth.Enabled {
			if auth == "" {
				logging.Log(logging.Error, "no authorization header found")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				logging.Log(logging.Error, "invalid authorization header")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				logging.Log(logging.Error, "failed to decode base64 token - %v", err)
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != alertSender.config.EasyAPIPromAlertSMS.Auth.Username || password != alertSender.config.EasyAPIPromAlertSMS.Auth.Password {
				logging.Log(logging.Error, "invalid username or password")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
