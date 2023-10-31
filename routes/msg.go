package routes

import (
	// "chat-be/controllers"

	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func msgRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/getMsgs", controllers.GetMsgs)

	return r
}
