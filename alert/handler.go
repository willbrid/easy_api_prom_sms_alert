package alert

import "net/http"

func (alertSender *AlertSender) AlertHandler(resp http.ResponseWriter, req *http.Request) {
	// Définir la logique d'envoie de SMS

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}
