package myMiddleware

import (
	"net/http"
	"os"
	"strings"
)

func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, rq *http.Request) {
		u, p, ok := rq.BasicAuth()
		if !ok || len(strings.TrimSpace(u)) < 1 || len(strings.TrimSpace(p)) < 1 {
			unauthorised(rw)
			return
		}

		if u != os.Getenv("DB_USER") || p != os.Getenv("DB_PASSWORD") {
			unauthorised(rw)
			return
		}

		handler(rw, rq)
	}
}

func unauthorised(rw http.ResponseWriter) {
	rw.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	rw.WriteHeader(http.StatusUnauthorized)
}
