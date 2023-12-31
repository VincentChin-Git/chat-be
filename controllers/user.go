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
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var ctx SignUpType
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	result, err := services.Signup(ctx.Mobile, ctx.Username, ctx.Password)

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

	_id, ok := r.Context().Value(middleware.ContextKey("parsedId")).(string)

	if !ok {
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	result, err := services.GetUserInfoByToken(_id)

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
		utils.JsonResponse(w, true, http.StatusOK)
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
		utils.JsonResponse(w, true, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

func SearchUser(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	token := ""
	var err error
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

	if len(authHeader) == 2 {
		token, err = utils.DecodeToken(authHeader[1])
		if err != nil {
			fmt.Println(err.Error())
			token = ""
		}
	}

	user, err := services.SearchUser(key, value, token)
	if err == nil {
		utils.JsonResponse(w, user, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

type ForgetPasswordType struct {
	Code     string `json:"code,omitempty"`
	Password string `json:"password,omitempty"`
	Id       string `json:"_id,omitempty"`
}

func ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var ctx ForgetPasswordType
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	err = services.ForgetPassword(ctx.Id, ctx.Code, ctx.Password)
	if err == nil {
		utils.JsonResponse(w, true, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

func AddForgetPassword(w http.ResponseWriter, r *http.Request) {
	var ctx struct {
		UserId string `json:"userId,omitempty"`
		Code   string `json:"code,omitempty"`
	}
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		fmt.Println(err)
		utils.JsonResponseError(w, "999999", "", http.StatusBadRequest)
		return
	}

	err = services.AddForgetPassword(ctx.UserId, ctx.Code)

	if err == nil {
		utils.JsonResponse(w, true, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}

}
func GetForgetPassCode(w http.ResponseWriter, r *http.Request) {
	mobile := r.URL.Query().Get("mobile")

	userId, err := services.GetForgetPassCode(mobile)

	if err == nil {
		utils.JsonResponse(w, userId, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}

}
