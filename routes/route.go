package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nexmedis-be-technical-test/handlers"
	v1 "github.com/nexmedis-be-technical-test/routes/v1"
)

func SetupRouter(handler handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger) // Optional: logging middleware

	r.Route("/v1", func(r chi.Router) {
		v1.AuthV1Routes(r, handler)
		v1.ProductV1Routes(r, handler)
		v1.CartV1Routes(r, handler)
		v1.TransactionV1Routes(r, handler)
	})

	return r
}
