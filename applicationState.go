package main

import (
	"sync"
	"time"
)

var StartTime = time.Now().Format("2006-01-02 15:04:05")
var mu sync.RWMutex
var GeoJsonPolygon = ""
var LastUpdatedAt = ""
