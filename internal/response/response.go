package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

type StatusResponse struct {
	Version       string `json:"version"`
	DeployFlag    string `json:"deploy_flag"`
	Polygon       string `json:"polygon"`
	Uptime        string `json:"uptime"`
	LastUpdatedAt string `json:"last_updated_at"`
}

func Error(message string, w http.ResponseWriter, httpCode int) {
	e := &ErrorResponse{false, message}
	jsonError, _ := json.Marshal(e)
	w.WriteHeader(httpCode)
	w.Write(jsonError)
}

func Success(payload map[string]interface{}, w http.ResponseWriter) {
	r := &SuccessResponse{Success: true, Data: payload}
	res, _ := json.Marshal(r)
	w.Write(res)
}
