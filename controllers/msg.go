package controllers

import (
	"chat-be/middleware"
	"chat-be/services"
	"chat-be/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetMsgs(w http.ResponseWriter, r *http.Request) {
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

	contactId := r.URL.Query().Get("contactId")
	if contactId == "" {
		utils.JsonResponseError(w, "999999", "Invalid Info", http.StatusBadRequest)
		return
	}

	result, err := services.GetMsgs(_id, contactId, skip, pageSize)
	if err != nil {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	} else {
		utils.JsonResponse(w, result, http.StatusOK)
	}
}

func UpdateMsgStatus(w http.ResponseWriter, r *http.Request) {
	var ctx []struct {
		Id     string `json:"_id,omitempty"`
		Status string `json:"status,omitempty"`
	}
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	for _, item := range ctx {
		err = services.UpdateMsgStatus(item.Id, item.Status, _id)
		if err != nil {
			utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
			return
		}
	}

	utils.JsonResponse(w, true, http.StatusOK)

}

func UpdateMsgToReceived(w http.ResponseWriter, r *http.Request) {
	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}
	err := services.UpdateMsgToReceived(_id)
	if err != nil {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	} else {
		utils.JsonResponse(w, true, http.StatusOK)
	}
}

func SendMsg(w http.ResponseWriter, r *http.Request) {
	var ctx struct {
		ReceiveId   string `json:"receiveId"`
		Content     string `json:"content,omitempty" bson:"content,omitempty"`
		ContentType string `json:"contentType,omitempty" bson:"contentType,omitempty"` /* text / image / video */
	}

	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	id, err := services.SendMsg(_id, ctx.ReceiveId, ctx.Content, ctx.ContentType)
	if err != nil {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	} else {
		utils.JsonResponse(w, id, http.StatusOK)
	}
}

func GetOverviewMsg(w http.ResponseWriter, r *http.Request) {
	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	search := r.URL.Query().Get("search")

	page, errPage := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, errPageSize := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if errPage != nil {
		page = 1
	}

	if errPageSize != nil {
		pageSize = 10
	}

	skip := utils.ToSkipRow(page, pageSize)

	result, err := services.GetOverviewMsg(_id, search, skip, pageSize)
	if err != nil {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	} else {
		utils.JsonResponse(w, result, http.StatusOK)
	}

}

func GetUnreadMsg(w http.ResponseWriter, r *http.Request) {
	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)
	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	result, err := services.GetUnreadMsg(_id)

	if err != nil {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	} else {
		utils.JsonResponse(w, result, http.StatusOK)
	}

}
