package gps

import (
	"math"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// Distance calculates the distance between two coordinates using the Haversine formula.
func Distance(c1, c2 Coordinate) float64 {
	const earthRadius = 6371 // in kilometers

	lat1 := toRadians(c1.Latitude)
	lat2 := toRadians(c2.Latitude)

	deltaLat := toRadians(c2.Latitude - c1.Latitude)
	deltaLon := toRadians(c2.Longitude - c1.Longitude)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c

	return distance
}

// toRadians converts a degree value to radians.
func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
