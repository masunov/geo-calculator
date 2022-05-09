package main

import (
	"encoding/json"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type StatusResponse struct {
	Version       string `json:"version"`
	DeployFlag    string `json:"deploy_flag"`
	Polygon       string `json:"polygon"`
	Uptime        string `json:"uptime"`
	LastUpdatedAt string `json:"last_updated_at"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	returnSuccessResponse(m, w)
	return
}

func statusHandler(w http.ResponseWriter, r *http.Request) {

	polygonStatus := "not exists"

	if len(strings.TrimSpace(GeoJsonPolygon)) != 0 {
		polygonStatus = "exists"
	}

	response := &StatusResponse{
		Version:       os.Getenv("APP_VERSION"),
		DeployFlag:    os.Getenv("APP_DEPLOY_FLAG"),
		Polygon:       polygonStatus,
		Uptime:        StartTime,
		LastUpdatedAt: LastUpdatedAt,
	}

	res, _ := json.Marshal(response)

	_, err := w.Write(res)
	if err != nil {
		return
	}

	return
}

func loadPolygonHandler(w http.ResponseWriter, r *http.Request) {

	GeoJsonPolygon = r.FormValue("polygon")

	if len(strings.TrimSpace(GeoJsonPolygon)) != 0 {
		LastUpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	return
}

func showPolygonHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["polygon"] = GeoJsonPolygon
	returnSuccessResponse(m, w)
	return
}

func checkPointHandler(w http.ResponseWriter, r *http.Request) {

	if len(strings.TrimSpace(GeoJsonPolygon)) == 0 {
		returnErrorResponse("Polygon not found", w, http.StatusBadRequest)
		return
	}

	featureCollection, _ := geojson.UnmarshalFeatureCollection([]byte(GeoJsonPolygon))

	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("lon"), 64)

	m := make(map[string]string)
	m["point_status"] = "out of polygon"

	// Pass in the feature collection + a point of Long/Lat
	if isPointInsidePolygon(featureCollection, orb.Point{lon, lat}) {
		m["point_status"] = "inside polygon"
	}

	returnSuccessResponse(m, w)
	return
}
