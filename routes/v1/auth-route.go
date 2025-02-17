package v1

import (
	"github.com/go-chi/chi"
	"github.com/nexmedis-be-technical-test/handlers"
)

func AuthV1Routes(r chi.Router, handler handlers.Handler) {
	r.Post("/register", handler.AuthRegisterUser)
	r.Post("/login", handler.AuthLogin)

	r.Route("/admin", func(r chi.Router) {
		r.Post("/register", handler.AuthRegisterAdmin)
	})
}
