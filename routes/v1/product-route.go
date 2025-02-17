package v1

import (
	"github.com/go-chi/chi"
	"github.com/nexmedis-be-technical-test/handlers"
)

func ProductV1Routes(r chi.Router, handler handlers.Handler) {
	r.Get("/products", handler.ProductGetListByFilter)
}
