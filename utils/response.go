package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	// You can include other fields like "status", "message", etc.
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Success bool        `json:"success"`
}

type ResponseError struct {
	ErrMessage string `json:"errMessage"`
	ErrCode    string `json:"errCode"`
}

func JsonResponse(w http.ResponseWriter, data interface{}, status int) {
	isSuccess := status == http.StatusOK
	response := Response{Data: data, Code: status, Success: isSuccess}
	responseJSON, errJSON := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, errWrite := w.Write(responseJSON)
	if errWrite != nil || errJSON != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}

func JsonResponseError(w http.ResponseWriter, errCode string, errMessage string) {
	realErr := "Please try again later."
	if errMessage != "" {
		realErr = errMessage
	}
	response := ResponseError{ErrCode: errCode, ErrMessage: realErr}
	responseJSON, errJSON := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, errWrite := w.Write(responseJSON)
	if errWrite != nil || errJSON != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}
