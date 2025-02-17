package v1

import (
	"github.com/go-chi/chi"
	"github.com/nexmedis-be-technical-test/handlers"
	"github.com/nexmedis-be-technical-test/middlewares"
)

func TransactionV1Routes(r chi.Router, handler handlers.Handler) {
	r.Route("/topup", func(r chi.Router) {
		// Apply the authentication middleware
		r.Use(middlewares.AuthenticateToken)
		r.Patch("/", handler.TransactionTopUpBalanceByUserId)
	})

	r.Route("/withdraw", func(r chi.Router) {
		r.Use(middlewares.AuthenticateToken)
		r.Patch("/", handler.TransactionWithdrawBalanceByUserId)
	})

	r.Route("/purchase", func(r chi.Router) {
		r.Use(middlewares.AuthenticateToken)
		r.Post("/", handler.TransactionPurchaseCart)
	})
}
