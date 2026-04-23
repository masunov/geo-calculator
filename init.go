package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/paulmach/orb/geojson"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		return
	}

	loadPolygonFromURL(os.Getenv("POLYGON_SOURCE_URL"))

	r := mux.NewRouter()

	r.Use(commonMiddleware)

	r.HandleFunc("/", indexHandler).Methods(http.MethodGet)
	r.HandleFunc("/status", statusHandler).Methods(http.MethodGet)
	r.HandleFunc("/load-polygon", loadPolygonHandler).Methods(http.MethodPost)
	r.HandleFunc("/show-polygon", showPolygonHandler).Methods(http.MethodGet)
	r.HandleFunc("/check-point", checkPointHandler).Methods(http.MethodGet)
	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), r)
	if err != nil {
		return
	}
}

func loadPolygonFromURL(url string) {
	if len(strings.TrimSpace(url)) == 0 {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch polygon from %s: %s", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read polygon response: %s", err)
		return
	}

	if _, err = geojson.UnmarshalFeatureCollection(body); err != nil {
		log.Printf("Invalid GeoJSON received from %s: %s", url, err)
		return
	}

	mu.Lock()
	GeoJsonPolygon = string(body)
	LastUpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	mu.Unlock()

	log.Printf("Polygon loaded from %s", url)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
