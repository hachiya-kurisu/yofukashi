// Package yofukashi implements functions useful for the nocturnal small web
package yofukashi

import (
	"math"
	"time"
)

const Version = "0.2.0"

// DawnDusk returns dawn and dusk on day t at latitude lat and longitude lon.
func DawnDusk(t time.Time, lat, lon float64) (time.Time, time.Time) {
	day := t.YearDay()
	x := math.Sin(360 * (float64(day) + 284) / 365.0 * math.Pi / 180)
	y := -math.Tan(lat*math.Pi/180) * math.Tan(23.44*x*math.Pi/180)
	if y < -1 {
		y = -1
	}
	if y > 1 {
		y = 1
	}
	hours := 1 / 15.0 * math.Acos(y) * 180 / math.Pi
	noon := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
	_, offsetSec := noon.Zone()
	utcoffset := float64(offsetSec) / 3600.0
	m := utcoffset * 15.0
	offset := (m - lon) / 15.0 * 3600
	noon = noon.Add(time.Duration(int64(offset)) * time.Second)
	d := time.Duration(int64(hours * 3600 * 1e9))
	return noon.Add(-d), noon.Add(d)
}

// Daytime returns true if t is after dawn and before dusk
func Daytime(t time.Time, lat, lon float64) bool {
	dawn, dusk := DawnDusk(t, lat, lon)
	return t.Before(dusk) && t.After(dawn)
}

// Nighttime returns true if t is after dusk and before dawn
func Nighttime(t time.Time, lat, lon float64) bool {
	return !Daytime(t, lat, lon)
}
