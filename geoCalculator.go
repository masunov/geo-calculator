package main

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

func isPointInsidePolygon(fc *geojson.FeatureCollection, point orb.Point) bool {
	for _, feature := range fc.Features {
		multiPoly, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			if planar.MultiPolygonContains(multiPoly, point) {
				return true
			}
		} else {
			polygon, isPoly := feature.Geometry.(orb.Polygon)
			if isPoly {
				if planar.PolygonContains(polygon, point) {
					return true
				}
			}
		}
	}
	return false
}
