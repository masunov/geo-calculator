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

type StatusResponse struct {
	Version       string `json:"version"`
	DeployFlag    string `json:"deploy_flag"`
	Polygon       string `json:"polygon"`
	Uptime        string `json:"uptime"`
	LastUpdatedAt string `json:"last_updated_at"`
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
