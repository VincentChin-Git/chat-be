package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func contactRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/getContact", controllers.GetContact)

	return r
}
