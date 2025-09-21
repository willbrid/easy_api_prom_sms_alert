package handler

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/handler/httphandler"
	"easy-api-prom-alert-sms/internal/handler/middleware"
	"easy-api-prom-alert-sms/internal/usecase"

	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Usecases *usecase.Usecases
	Router   *mux.Router
}

func NewHandler(usecases *usecase.Usecases, router *mux.Router) *Handler {
	return &Handler{usecases, router}
}

func (h *Handler) InitRouter(cfg *config.Config) {
	h.Router.Use(func(h http.Handler) http.Handler {
		return middleware.AuthMiddleware(h, cfg)
	})

	httphandler := httphandler.NewHTTPHandler(h.Usecases)

	h.Router.HandleFunc("/healthz", httphandler.HandleHealthCheck).Methods("GET")
	h.Router.HandleFunc("/api-alert", httphandler.HandleAlert).Methods("POST")
}
