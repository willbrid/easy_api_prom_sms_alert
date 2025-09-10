package handler

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/middleware"
	"easy-api-prom-alert-sms/internal/usecase"
	"easy-api-prom-alert-sms/pkg/logger"

	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(router *mux.Router, cfg *config.Config, a usecase.IAlert, l logger.ILogger) {
	handler := NewHandler(a, l)

	router.Use(func(h http.Handler) http.Handler {
		return middleware.AuthMiddleware(h, cfg)
	})

	router.HandleFunc("/healthz", handler.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/api-alert", handler.HandleAlert).Methods("POST")
}
