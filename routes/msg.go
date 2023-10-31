package routes

import (
	// "chat-be/controllers"

	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func msgRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/getMsgs", controllers.GetMsgs)
	r.Post("/updateMsgStatus", controllers.UpdateMsgStatus)
	r.Post("/sendMsg", controllers.SendMsg)
	r.Get("/getOverviewMsg", controllers.GetOverviewMsg)

	return r
}
