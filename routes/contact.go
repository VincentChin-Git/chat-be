package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func contactRoutes() chi.Router {
	r := chi.NewRouter()

	// get
	r.Get("/getContacts", controllers.GetContacts)
	r.Get("/getContact", controllers.GetContact)

	// post
	r.Post("/addContact", controllers.AddContact)
	r.Post("/removeContact", controllers.RemoveContact)

	return r
}
