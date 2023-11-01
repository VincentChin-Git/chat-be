package routes

import (
	"github.com/go-chi/chi/v5"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()

	r.Mount("/aws", awsRoutes())
	r.Mount("/contact", contactRoutes())
	r.Mount("/msg", msgRoutes())
	r.Mount("/user", userRoutes())

	return r
}
