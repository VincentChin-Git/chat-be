package routes

import (
	"github.com/go-chi/chi/v5"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()

	r.Mount("/user", userRoutes())
	r.Mount("/contact", contactRoutes())

	return r
}
