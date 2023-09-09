package routes

import (
	"github.com/FranMT-S/Challenge-Go/src/controller"
	"github.com/go-chi/chi/v5"
)

func MailRouter() *chi.Mux {
	router := chi.NewMux()

	router.Route(MAIL_ROUTE, func(r chi.Router) {
		r.Get("/", controller.GetTotalMessage)
	})

	return router
}
