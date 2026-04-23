package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	"GeoCalc/internal/geo"
	"GeoCalc/internal/response"
	"GeoCalc/internal/state"
)

func Index(w http.ResponseWriter, r *http.Request) {
	response.Success(make(map[string]interface{}), w)
}

func Status(w http.ResponseWriter, r *http.Request) {
	polygon := state.GetPolygon()

	polygonStatus := "not exists"
	if len(strings.TrimSpace(polygon)) != 0 {
		polygonStatus = "exists"
	}

	sr := &response.StatusResponse{
		Version:       os.Getenv("APP_VERSION"),
		DeployFlag:    os.Getenv("APP_DEPLOY_FLAG"),
		Polygon:       polygonStatus,
		Uptime:        state.StartTime(),
		LastUpdatedAt: state.GetLastUpdatedAt(),
	}

	res, _ := json.Marshal(sr)
	w.Write(res)
}

func LoadPolygon(w http.ResponseWriter, r *http.Request) {
	polygon := r.FormValue("polygon")

	if len(strings.TrimSpace(polygon)) == 0 {
		response.Error("Polygon is required", w, http.StatusBadRequest)
		return
	}

	if _, err := geojson.UnmarshalFeatureCollection([]byte(polygon)); err != nil {
		response.Error("Invalid GeoJSON", w, http.StatusBadRequest)
		return
	}

	state.SetPolygon(polygon)
	response.Success(make(map[string]interface{}), w)
}

func ShowPolygon(w http.ResponseWriter, r *http.Request) {
	polygon := state.GetPolygon()

	m := make(map[string]interface{})
	if len(strings.TrimSpace(polygon)) == 0 {
		m["polygon"] = nil
	} else {
		m["polygon"] = polygon
	}
	response.Success(m, w)
}

func CheckPoint(w http.ResponseWriter, r *http.Request) {
	polygon := state.GetPolygon()

	if len(strings.TrimSpace(polygon)) == 0 {
		response.Error("Polygon not found", w, http.StatusBadRequest)
		return
	}

	featureCollection, err := geojson.UnmarshalFeatureCollection([]byte(polygon))
	if err != nil {
		response.Error("Invalid polygon GeoJSON", w, http.StatusInternalServerError)
		return
	}

	lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	if err != nil {
		response.Error("Invalid lat value", w, http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(r.FormValue("lon"), 64)
	if err != nil {
		response.Error("Invalid lon value", w, http.StatusBadRequest)
		return
	}

	m := make(map[string]interface{})
	m["point_status"] = "out of polygon"

	if geo.IsPointInsidePolygon(featureCollection, orb.Point{lon, lat}) {
		m["point_status"] = "inside polygon"
	}

	response.Success(m, w)
}
