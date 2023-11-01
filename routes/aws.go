package routes

import (
	"chat-be/controllers"

	"github.com/go-chi/chi/v5"
)

func awsRoutes() chi.Router {
	r := chi.NewRouter()

	// get

	// post
	r.Get("/uploadImgSignature", controllers.UploadImgSignature)
	r.Get("/uploadVideoSignature", controllers.UploadVideoSignature)

	return r
}
