package gps

import (
	"math"
	"testing"
)

func TestDistanceRoundedAntipodes(t *testing.T) {
	pairs := [][2]Coordinate{
		{{Latitude: -88.5, Longitude: 2.5}, {Latitude: 88.5, Longitude: -177.5}},
		{{Latitude: 51.5074, Longitude: -0.1278}, {Latitude: -51.5074, Longitude: 179.8722}},
		{{Latitude: -90, Longitude: 0}, {Latitude: 90, Longitude: 0}},
	}
	want := math.Pi * 6371
	for _, pair := range pairs {
		for _, points := range [][2]Coordinate{pair, {pair[1], pair[0]}} {
			got := Distance(points[0], points[1])
			if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-want) > 0.001 {
				t.Errorf("Distance(%+v, %+v) = %v, want approximately %v", points[0], points[1], got, want)
			}
		}
	}
}

func TestDistanceAcrossDateLine(t *testing.T) {
	a := Coordinate{Latitude: 0, Longitude: 179}
	b := Coordinate{Latitude: 0, Longitude: -179}
	want := 6371 * 2 * math.Pi / 180
	if got := Distance(a, b); math.Abs(got-want) > 1e-9 {
		t.Errorf("Distance across date line = %v, want %v", got, want)
	}
	if Distance(a, b) != Distance(b, a) {
		t.Error("Distance is not symmetric")
	}
}
