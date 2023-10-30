package controllers

import (
	"chat-be/services"
	"chat-be/utils"
	"encoding/json"
	"fmt"
	"net/http"
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
		utils.JsonResponseError(w, "999999", "")
		return
	}

	result, err := services.Signup(ctx.Mobile, ctx.Username, ctx.Password, ctx.Nickname)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error())
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
		utils.JsonResponseError(w, "999999", "")
		return
	}

	result, err := services.Login(ctx.Param, ctx.Password)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error())
	}
}

type GetUserInfoByTokenType struct {
	Token string `json:"token"`
}

func GetUserInfoByToken(w http.ResponseWriter, r *http.Request) {
	var ctx GetUserInfoByTokenType
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "")
		return
	}

	result, err := services.GetUserInfoByToken(ctx.Token)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error())
	}
}
