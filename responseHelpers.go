package main

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
}

func returnErrorResponse(errorMessage string, w http.ResponseWriter, httpCode int) {
	e := &ErrorResponse{
		false,
		errorMessage,
	}
	jsonError, _ := json.Marshal(e)
	w.WriteHeader(httpCode)
	_, err := w.Write(jsonError)
	if err != nil {
		return
	}
}

func returnSuccessResponse(payload map[string]string, w http.ResponseWriter) {
	response := &SuccessResponse{
		Success: true,
		Data:    payload,
	}
	res, _ := json.Marshal(response)
	_, err := w.Write(res)
	if err != nil {
		return
	}
}
