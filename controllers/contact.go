package controllers

import (
	"chat-be/middleware"
	"chat-be/services"
	"chat-be/utils"
	"net/http"
	"strconv"
)

func GetContact(w http.ResponseWriter, r *http.Request) {
	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	page, errPage := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, errPageSize := strconv.Atoi(r.URL.Query().Get("pageSize"))

	skip := utils.ToSkipRow(page, pageSize)

	if errPage != nil || errPageSize != nil {
		utils.JsonResponseError(w, "999999", "Invalid Info", http.StatusBadRequest)
	}

	result, err := services.GetContact(_id, skip, pageSize)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}
