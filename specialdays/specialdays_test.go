package specialdays

import (
	"testing"
	"time"
)

// TestCalculateEasterSunday testet die Berechnung des Ostersonntags für verschiedene Jahre.
func TestCalculateEasterSunday(t *testing.T) {
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2020, time.Date(2020, time.April, 12, 0, 0, 0, 0, time.UTC)},
		{2024, time.Date(2024, time.March, 31, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, time.April, 20, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, time.April, 5, 0, 0, 0, 0, time.UTC)},
		// Weitere Testfälle...
	}

	for _, tc := range testCases {
		calculatedDate := CalculateEasterSunday(tc.year)
		if !tc.expected.Equal(calculatedDate) {
			t.Errorf("Expected Easter Sunday of %d to be %v, but got %v", tc.year, tc.expected, calculatedDate)
		}
	}
}

// TestCalculateMothersDay testet die Berechnung des Muttertags für verschiedene Jahre.
func TestCalculateMothersDay(t *testing.T) {
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2020, time.Date(2020, time.May, 10, 0, 0, 0, 0, time.UTC)},
		{2024, time.Date(2024, time.May, 12, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, time.May, 11, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, time.May, 10, 0, 0, 0, 0, time.UTC)},
		// Weitere Testfälle...
	}

	for _, tc := range testCases {
		calculatedDate := CalculateMothersDay(tc.year)
		if !tc.expected.Equal(calculatedDate) {
			t.Errorf("Expected Mothers Day of %d to be %v, but got %v", tc.year, tc.expected, calculatedDate)
		}
	}
}

// TestGetSpecialDays testet die Map der besonderen Tage für ein bestimmtes Jahr.
func TestGetSpecialDays(t *testing.T) {
	year := 2024
	specialDays := GetSpecialDays(year)

	testCases := []struct {
		name     string
		expected time.Time
	}{
		{"Neujahrstag", time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"Tag der Arbeit", time.Date(year, time.May, 1, 0, 0, 0, 0, time.UTC)},
		{"Tag der Deutschen Einheit", time.Date(year, time.October, 3, 0, 0, 0, 0, time.UTC)},
		{"Heiligabend", time.Date(year, time.December, 24, 0, 0, 0, 0, time.UTC)},
		// Weitere Testfälle...
	}

	for _, tc := range testCases {
		if date, ok := specialDays[tc.name]; ok {
			if !tc.expected.Equal(date) {
				t.Errorf("Expected %s to be %v, but got %v", tc.name, tc.expected, date)
			}
		} else {
			t.Errorf("Expected %s to be in special days map", tc.name)
		}
	}
}
