package views

import (
	"database/sql"
	"net/http"
)

type View struct {
	Route       string
	Middlewares []func(http.Handler) http.Handler
	Handler     http.Handler
	Methods     []string
}

var DB *sql.DB

// Middleware chaining
func chainMiddlewares(handler http.Handler, middlewares []func(http.Handler) http.Handler) http.Handler {
	for i := 0; i < len(middlewares); i++ { // Apply middlewares in normal order
		handler = middlewares[i](handler)
	}
	return handler
}

// Routes function
func Routes(mux *http.ServeMux, views []View) {
	for _, view := range views {
		handlerWithMiddlewares := chainMiddlewares(view.Handler, view.Middlewares)
		for _, method := range view.Methods {
			mux.HandleFunc(view.Route, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
					return
				}
				handlerWithMiddlewares.ServeHTTP(w, r)
			})
		}
	}
}
