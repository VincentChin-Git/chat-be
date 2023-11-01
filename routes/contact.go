package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func contactRoutes() chi.Router {
	r := chi.NewRouter()

	// get
	r.Get("/getContact", controllers.GetContact)

	// post
	r.Post("/addContact", controllers.AddContact)

	return r
}
