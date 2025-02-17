package v1

import (
	"github.com/go-chi/chi"
	"github.com/nexmedis-be-technical-test/handlers"
	"github.com/nexmedis-be-technical-test/middlewares"
)

func CartV1Routes(r chi.Router, handler handlers.Handler) {
	r.Route("/carts", func(r chi.Router) {
		// Apply the authentication middleware
		r.Use(middlewares.AuthenticateToken)
		r.Post("/", handler.UserCartAddItem)
	})
}
