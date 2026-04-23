package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/paulmach/orb/geojson"

	"GeoCalc/internal/handlers"
	"GeoCalc/internal/state"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	loadPolygonFromURL(os.Getenv("POLYGON_SOURCE_URL"))

	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.HandleFunc("/", handlers.Index).Methods(http.MethodGet)
	r.HandleFunc("/status", handlers.Status).Methods(http.MethodGet)
	r.HandleFunc("/load-polygon", handlers.LoadPolygon).Methods(http.MethodPost)
	r.HandleFunc("/show-polygon", handlers.ShowPolygon).Methods(http.MethodGet)
	r.HandleFunc("/check-point", handlers.CheckPoint).Methods(http.MethodGet)

	if err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), r); err != nil {
		log.Fatal(err)
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

	state.SetPolygon(string(body))
	log.Printf("Polygon loaded from %s", url)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
