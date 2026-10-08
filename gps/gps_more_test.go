package gps

import "testing"

func TestDistanceZeroAndSame(t *testing.T) {
	a := Coordinate{Latitude: 0, Longitude: 0}
	if d := Distance(a, a); d != 0 {
		t.Fatalf("expected zero distance, got %v", d)
	}
}

func TestDistanceAntipode(t *testing.T) {
	// Approximate antipodal points (opposite sides of earth)
	a := Coordinate{Latitude: 0, Longitude: 0}
	b := Coordinate{Latitude: 0, Longitude: 180}
	d := Distance(a, b)
	// Earth's half-circumference ~ 20037.5 km
	if d < 20000 || d > 20100 {
		t.Fatalf("antipode distance out of expected range: %v", d)
	}
}
