package myServer

import (
	"fmt"
	"net/http"

	myMiddleware "github.com/FranMT-S/Challenge-Go/src/middleware"
	"github.com/FranMT-S/Challenge-Go/src/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Iserve interface {
	Start()
}

var _server Iserve

type server struct {
	router *chi.Mux
}

func Server() Iserve {

	if _server == nil {
		_server = &server{router: chi.NewRouter()}
	}

	return _server
}

func (_server *server) Start() {

	_server.router.Use(middleware.Logger)
	_server.MountHandlers()

	fmt.Printf("Servidor escuchando en el puerto %s...\n", "3000")
	http.ListenAndServe(":3000", _server.router)

}

func (s *server) MountHandlers() {
	// Mount all Middleware here
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.AllowContentType("application/json"))
	s.router.Use(myMiddleware.JsonMiddleware)
	// s.router.Use(myMiddleware.ZincHeader)

	// Mount all handlers here
	// s.router.Get("/", HelloWorld)
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/", func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Hello World!"))
		})

		r.Mount("/", routes.MailRouter())
	})

}
