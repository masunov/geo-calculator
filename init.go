package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		return
	}

	r := mux.NewRouter()

	r.Use(commonMiddleware)

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/status", statusHandler)
	r.HandleFunc("/load-polygon", loadPolygonHandler)
	r.HandleFunc("/show-polygon", showPolygonHandler)
	r.HandleFunc("/check-point", checkPointHandler)
	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), r)
	if err != nil {
		return
	}
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
