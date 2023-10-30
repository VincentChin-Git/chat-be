package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func testRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {})

	return r
}
