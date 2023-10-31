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

	skip := utils.ToSkipRow(page, pageSize)

	if errPage != nil || errPageSize != nil {
		utils.JsonResponseError(w, "999999", "Invalid Info", http.StatusBadRequest)
		return
	}

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

func updateMsgStatus(w http.ResponseWriter, r *http.Request) {
	var ctx struct {
		Id     string `json:"_id,omitempty"`
		Status string `json:"status,omitempty"`
	}
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	err = services.UpdateMsgStatus(ctx.Id, ctx.Status)
	if err != nil {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	} else {
		utils.JsonResponse(w, "", http.StatusOK)

	}

}
