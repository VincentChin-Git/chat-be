package controllers

import (
	"chat-be/services"
	"chat-be/utils"
	"net/http"
)

func UploadImgSignature(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	result, err := services.UploadImgSignature(filename)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}

func UploadVideoSignature(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	result, err := services.UploadVideoSignature(filename)

	if err == nil {
		utils.JsonResponse(w, result, http.StatusOK)
	} else {
		utils.JsonResponseError(w, "999999", err.Error(), http.StatusBadRequest)
	}
}
