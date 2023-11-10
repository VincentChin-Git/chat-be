package routes

import (
	// "chat-be/controllers"

	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func msgRoutes() chi.Router {
	r := chi.NewRouter()

	// get
	r.Get("/getMsgs", controllers.GetMsgs)
	r.Get("/getOverviewMsg", controllers.GetOverviewMsg)
	r.Get("/getUnreadMsg", controllers.GetUnreadMsg)

	// post
	r.Post("/sendMsg", controllers.SendMsg)
	r.Post("/updateMsgStatus", controllers.UpdateMsgStatus)
	r.Post("/updateMsgToReceived", controllers.UpdateMsgToReceived)

	return r
}
