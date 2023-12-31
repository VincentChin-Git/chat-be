package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func userRoutes() chi.Router {
	r := chi.NewRouter()

	// get
	r.Get("/getForgetPassCode", controllers.GetForgetPassCode)
	r.Get("/getUserInfoByToken", controllers.GetUserInfoByToken)
	r.Get("/searchUser", controllers.SearchUser)

	// post
	r.Post("/addForgetPassword", controllers.AddForgetPassword)
	r.Post("/changePassword", controllers.ChangePassword)
	r.Post("/forgetPassword", controllers.ForgetPassword)
	r.Post("/login", controllers.Login)
	r.Post("/signup", controllers.Signup)
	r.Post("/updateUserInfo", controllers.UpdateUserInfo)

	return r
}
