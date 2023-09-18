package config

import (
	"net/http"

	"github.com/mercadolibre/fury_go-core/pkg/web"
)

// JSONResponse adds the content-type json to every request.
func JSONResponse() web.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler(w, r)
		}
	}
}
