package httpapi

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"crypto-go/internal/httpapi/openapi"
)

func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/crypto/encrypt", handler.Encrypt)
	mux.HandleFunc("POST /api/crypto/decrypt", handler.Decrypt)
	mux.HandleFunc("POST /api/crypto/hash", handler.Hash)
	mux.HandleFunc("GET /api/crypto/logs", handler.Logs)
	mux.HandleFunc("GET /api/crypto/logs/{id}", handler.LogByID)

	mux.HandleFunc("GET /v3/api-docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(openapi.Spec)
	})
	mux.Handle("GET /swagger-ui/", httpSwagger.Handler(
		httpSwagger.URL("/v3/api-docs"),
	))

	return mux
}
