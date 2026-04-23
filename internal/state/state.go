package state

import (
	"sync"
	"time"
)

var startTime = time.Now().Format("2006-01-02 15:04:05")
var mu sync.RWMutex
var geoJsonPolygon = ""
var lastUpdatedAt = ""

func StartTime() string {
	return startTime
}

func GetPolygon() string {
	mu.RLock()
	defer mu.RUnlock()
	return geoJsonPolygon
}

func SetPolygon(polygon string) {
	mu.Lock()
	defer mu.Unlock()
	geoJsonPolygon = polygon
	lastUpdatedAt = time.Now().Format("2006-01-02 15:04:05")
}

func GetLastUpdatedAt() string {
	mu.RLock()
	defer mu.RUnlock()
	return lastUpdatedAt
}
