package utils

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(string(responseJSON), "returned Data")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, errWrite := w.Write(responseJSON)
	if errWrite != nil || errJSON != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}

func JsonResponseError(w http.ResponseWriter, errCode string, errMessage string, status int) {
	var realErr string
	if status == http.StatusBadRequest {
		realErr = "Please try again later"
	} else if status == http.StatusUnauthorized {
		realErr = "User unauthorized"
	}

	if errMessage != "" {
		realErr = errMessage
	}
	response := ResponseError{ErrCode: errCode, ErrMessage: realErr}
	responseJSON, errJSON := json.Marshal(response)
	fmt.Println(string(responseJSON), "returned Error")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, errWrite := w.Write(responseJSON)
	if errWrite != nil || errJSON != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}
