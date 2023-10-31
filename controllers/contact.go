package controllers

import (
	"chat-be/middleware"
	"chat-be/services"
	"chat-be/utils"
	"net/http"
)

func GetContact(w http.ResponseWriter, r *http.Request) {
	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	result, err := services.GetContact(_id)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}
