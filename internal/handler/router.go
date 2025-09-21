package handler

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/handler/httphandler"
	"easy-api-prom-alert-sms/internal/handler/middleware"
	"easy-api-prom-alert-sms/internal/usecase"

	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(router *mux.Router, cfg *config.Config, a usecase.IAlert) {
	handler := httphandler.NewHandler(a)

	router.Use(func(h http.Handler) http.Handler {
		return middleware.AuthMiddleware(h, cfg)
	})

	router.HandleFunc("/healthz", handler.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/api-alert", handler.HandleAlert).Methods("POST")
}
