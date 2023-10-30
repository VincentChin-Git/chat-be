package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func userRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/signup", controllers.Signup)
	r.Get("/login", controllers.Login)
	r.Get("/getUserInfoByToken", controllers.GetUserInfoByToken)

	return r
}
