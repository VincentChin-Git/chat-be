package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func awsRoutes() chi.Router {
	r := chi.NewRouter()

	// get

	// post
	r.Post("/UploadImgSignature", controllers.UploadImgSignature)
	r.Post("/UploadVideoSignature", controllers.UploadVideoSignature)

	return r
}
