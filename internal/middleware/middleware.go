package middleware

import (
	"easy-api-prom-alert-sms/config"

	"encoding/base64"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var auth string = req.Header.Get("Authorization")

		if cfg.EasyAPIPromAlertSMS.Auth.Enabled {
			if auth == "" {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != cfg.EasyAPIPromAlertSMS.Auth.Username || password != cfg.EasyAPIPromAlertSMS.Auth.Password {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
