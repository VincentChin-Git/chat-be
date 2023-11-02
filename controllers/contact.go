package controllers

import (
	"chat-be/middleware"
	"chat-be/services"
	"chat-be/utils"
	"encoding/json"
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

	if errPage != nil {
		page = 1
	}

	if errPageSize != nil {
		pageSize = 10
	}

	skip := utils.ToSkipRow(page, pageSize)

	result, err := services.GetContact(_id, skip, pageSize)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

func AddContact(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	var ctx struct {
		ContactId string `json:"contactId"`
	}
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	result, err := services.AddContact(userId, ctx.ContactId)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}

}

func RemoveContact(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	var ctx struct {
		ContactId string `json:"contactId"`
	}
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	err = services.RemoveContact(userId, ctx.ContactId)

	if err == nil {
		utils.JsonResponse(w, true, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}
