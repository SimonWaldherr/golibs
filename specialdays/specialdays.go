package specialdays

import (
	"time"
)

// CalculateEasterSunday berechnet den Ostersonntag für ein gegebenes Jahr.
func CalculateEasterSunday(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// CalculateMothersDay berechnet den Muttertag für ein gegebenes Jahr (2. Sonntag im Mai).
func CalculateMothersDay(year int) time.Time {
	// Der 1. Mai des Jahres
	mayFirst := time.Date(year, time.May, 1, 0, 0, 0, 0, time.UTC)

	// Wochentag des 1. Mai berechnen
	weekday := mayFirst.Weekday()

	// Den ersten Sonntag im Mai finden
	daysUntilSunday := (7 - int(weekday)) % 7
	firstSunday := mayFirst.AddDate(0, 0, daysUntilSunday)

	// Den zweiten Sonntag im Mai berechnen
	secondSunday := firstSunday.AddDate(0, 0, 7)

	return secondSunday
}

// CalculateFronleichnam berechnet Fronleichnam für ein gegebenes Jahr (60 Tage nach Ostern).
func CalculateFronleichnam(year int) time.Time {
	easterSunday := CalculateEasterSunday(year)
	return easterSunday.AddDate(0, 0, 60)
}

// CalculateFirstAdvent berechnet den ersten Advent für ein gegebenes Jahr.
func CalculateFirstAdvent(year int) time.Time {
	// Der erste Advent ist der Sonntag zwischen dem 27. November und dem 3. Dezember.
	firstAdvent := time.Date(year, time.December, 3, 0, 0, 0, 0, time.UTC)
	for firstAdvent.Weekday() != time.Sunday {
		firstAdvent = firstAdvent.AddDate(0, 0, -1)
	}
	return firstAdvent
}

// CalculateBussUndBettag berechnet den Buß- und Bettag für ein gegebenes Jahr.
func CalculateBussUndBettag(year int) time.Time {
	// Der Buß- und Bettag ist der Mittwoch vor dem ersten Advent.
	firstAdvent := CalculateFirstAdvent(year)
	bussUndBettag := firstAdvent.AddDate(0, 0, -11)
	
	return bussUndBettag
}

// GetSpecialDays gibt eine Map mit allen besonderen Tagen für ein gegebenes Jahr zurück.
func GetSpecialDays(year int) map[string]time.Time {
	specialDays := make(map[string]time.Time)

	// Statische Feiertage und besondere Tage
	specialDays["Neujahrstag"] = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	specialDays["Tag der Arbeit"] = time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	specialDays["Tag der Deutschen Einheit"] = time.Date(year, 10, 3, 0, 0, 0, 0, time.UTC)
	specialDays["Nikolaustag"] = time.Date(year, time.December, 6, 0, 0, 0, 0, time.UTC)
	specialDays["Heiligabend"] = time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC)
	specialDays["1. Weihnachtstag"] = time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	specialDays["2. Weihnachtstag"] = time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	specialDays["Silvester"] = time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)

	// Bewegliche Feiertage
	easterSunday := CalculateEasterSunday(year)
	specialDays["Karfreitag"] = easterSunday.AddDate(0, 0, -2)
	specialDays["Ostersonntag"] = easterSunday
	specialDays["Ostermontag"] = easterSunday.AddDate(0, 0, 1)
	specialDays["Christi Himmelfahrt"] = easterSunday.AddDate(0, 0, 39)
	specialDays["Pfingstsonntag"] = easterSunday.AddDate(0, 0, 49)
	specialDays["Pfingstmontag"] = easterSunday.AddDate(0, 0, 50)
	specialDays["Fronleichnam"] = easterSunday.AddDate(0, 0, 60)
	specialDays["erster Advent"] = CalculateFirstAdvent(year)

	// Andere besondere Tage
	specialDays["Valentinstag"] = time.Date(year, 2, 14, 0, 0, 0, 0, time.UTC)
	specialDays["Muttertag"] = CalculateMothersDay(year)
	specialDays["Halloween"] = time.Date(year, 10, 31, 0, 0, 0, 0, time.UTC)
	specialDays["Karnevalsbeginn"] = time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)

	// Weitere Feiertage
	specialDays["Heilige Drei Könige"] = time.Date(year, time.January, 6, 0, 0, 0, 0, time.UTC)
	specialDays["Fronleichnam"] = CalculateFronleichnam(year)
	specialDays["Mariä Himmelfahrt"] = time.Date(year, time.August, 15, 0, 0, 0, 0, time.UTC)
	specialDays["Allerheiligen"] = time.Date(year, time.November, 1, 0, 0, 0, 0, time.UTC)
	specialDays["Buß- und Bettag"] = CalculateBussUndBettag(year)

	return specialDays
}
