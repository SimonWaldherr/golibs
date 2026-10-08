// Package gps contains helpers for geographic coordinates.
package gps

import (
	"math"
)

// Coordinate represents a geographic point (latitude, longitude).
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// Distance returns the great-circle distance in kilometers between two
// geographic coordinates using the Haversine formula. The result is an
// approximation suitable for many applications.
func Distance(c1, c2 Coordinate) float64 {
	const earthRadius = 6371 // in kilometers

	lat1 := toRadians(c1.Latitude)
	lat2 := toRadians(c2.Latitude)

	deltaLat := toRadians(c2.Latitude - c1.Latitude)
	deltaLon := toRadians(c2.Longitude - c1.Longitude)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	// Roundoff can put a outside [0, 1] near antipodal coordinates.
	a = math.Max(0, math.Min(1, a))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c

	return distance
}

// toRadians converts degrees to radians.
func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
