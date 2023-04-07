package gps

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	tests := []struct {
		c1       Coordinate
		c2       Coordinate
		expected float64
	}{
		{
			Coordinate{Latitude: 51.5074, Longitude: -0.1278},  // London, UK
			Coordinate{Latitude: 40.7128, Longitude: -74.0060}, // New York City, USA
			5570.22,
		},
		{
			Coordinate{Latitude: 37.7749, Longitude: -122.4194}, // San Francisco, USA
			Coordinate{Latitude: 48.8566, Longitude: 2.3522},    // Paris, France
			8953.39,
		},
		{
			Coordinate{Latitude: 40.4168, Longitude: -3.7038},  // Madrid, Spain
			Coordinate{Latitude: 31.2304, Longitude: 121.4737}, // Shanghai, China
			10255.48,
		},
	}

	for _, test := range tests {
		distance := Distance(test.c1, test.c2)

		if math.Abs(distance-test.expected) > 0.1 {
			t.Errorf("Distance between %v and %v was %.2f km, expected %.2f km", test.c1, test.c2, distance, test.expected)
		}
	}
}
