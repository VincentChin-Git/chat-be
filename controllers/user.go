package controllers

import (
	"chat-be/middleware"
	"chat-be/models"
	"chat-be/services"
	"chat-be/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SignUpType struct {
	Mobile   string `json:"mobile"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var ctx SignUpType
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	result, err := services.Signup(ctx.Mobile, ctx.Username, ctx.Password, ctx.Nickname)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

type LoginType struct {
	Param    string `json:"param"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var ctx LoginType
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	result, err := services.Login(ctx.Param, ctx.Password)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

func GetUserInfoByToken(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	token := strings.Split(reqToken, "Bearer ")[1]

	result, err := services.GetUserInfoByToken(token)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User
	err := json.NewDecoder(r.Body).Decode(&userInfo)
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
	err = services.UpdateUserInfo(_id, userInfo)

	if err == nil {
		utils.JsonResponse(w, "", http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

type ChangePasswordType struct {
	OldPass string `json:"oldPass"`
	NewPass string `json:"newPass"`
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var ctx ChangePasswordType
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

	err = services.ChangePassword(_id, ctx.OldPass, ctx.NewPass)

	if err == nil {
		utils.JsonResponse(w, "", http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

func SearchUser(w http.ResponseWriter, r *http.Request) {
	var ctx struct {
		Mobile string `json:"mobile,omitempty"`
	}

	user, err := services.SearchUser(ctx.Mobile)
	if err == nil {
		utils.JsonResponse(w, user, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}
