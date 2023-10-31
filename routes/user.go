package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func userRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/signup", controllers.Signup)
	r.Post("/login", controllers.Login)
	r.Get("/getUserInfoByToken", controllers.GetUserInfoByToken)
	r.Post("/updateUserInfo", controllers.UpdateUserInfo)
	r.Post("/changePassword", controllers.ChangePassword)
	r.Get("/getContact", controllers.GetContact)

	return r
}
