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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	returnSuccessResponse(m, w)
	return
}

func statusHandler(w http.ResponseWriter, r *http.Request) {

	mu.RLock()
	polygon := GeoJsonPolygon
	lastUpdated := LastUpdatedAt
	mu.RUnlock()

	polygonStatus := "not exists"
	if len(strings.TrimSpace(polygon)) != 0 {
		polygonStatus = "exists"
	}

	response := &StatusResponse{
		Version:       os.Getenv("APP_VERSION"),
		DeployFlag:    os.Getenv("APP_DEPLOY_FLAG"),
		Polygon:       polygonStatus,
		Uptime:        StartTime,
		LastUpdatedAt: lastUpdated,
	}

	res, _ := json.Marshal(response)

	_, err := w.Write(res)
	if err != nil {
		return
	}

	return
}

func loadPolygonHandler(w http.ResponseWriter, r *http.Request) {

	polygon := r.FormValue("polygon")

	if len(strings.TrimSpace(polygon)) == 0 {
		returnErrorResponse("Polygon is required", w, http.StatusBadRequest)
		return
	}

	if _, err := geojson.UnmarshalFeatureCollection([]byte(polygon)); err != nil {
		returnErrorResponse("Invalid GeoJSON", w, http.StatusBadRequest)
		return
	}

	mu.Lock()
	GeoJsonPolygon = polygon
	LastUpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	mu.Unlock()

	m := make(map[string]interface{})
	returnSuccessResponse(m, w)

	return
}

func showPolygonHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	polygon := GeoJsonPolygon
	mu.RUnlock()

	m := make(map[string]interface{})
	if len(strings.TrimSpace(polygon)) == 0 {
		m["polygon"] = nil
	} else {
		m["polygon"] = polygon
	}
	returnSuccessResponse(m, w)
	return
}

func checkPointHandler(w http.ResponseWriter, r *http.Request) {

	mu.RLock()
	polygon := GeoJsonPolygon
	mu.RUnlock()

	if len(strings.TrimSpace(polygon)) == 0 {
		returnErrorResponse("Polygon not found", w, http.StatusBadRequest)
		return
	}

	featureCollection, err := geojson.UnmarshalFeatureCollection([]byte(polygon))
	if err != nil {
		returnErrorResponse("Invalid polygon GeoJSON", w, http.StatusInternalServerError)
		return
	}

	lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	if err != nil {
		returnErrorResponse("Invalid lat value", w, http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(r.FormValue("lon"), 64)
	if err != nil {
		returnErrorResponse("Invalid lon value", w, http.StatusBadRequest)
		return
	}

	m := make(map[string]interface{})
	m["point_status"] = "out of polygon"

	if isPointInsidePolygon(featureCollection, orb.Point{lon, lat}) {
		m["point_status"] = "inside polygon"
	}

	returnSuccessResponse(m, w)
	return
}
