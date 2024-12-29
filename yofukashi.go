package yofukashi

import (
	"math"
	"time"
)

const Version = "0.0.4"

func DawnDusk(t time.Time, lat float64) (time.Time, time.Time) {
	day := t.YearDay()
	x := math.Sin(360 * (float64(day) + 284) / 365.0 * math.Pi / 180)
	y := -math.Tan(lat*math.Pi/180) * math.Tan(23.44*x*math.Pi/180)
	hours := 1 / 15.0 * math.Acos(y) * 180 / math.Pi
	noon := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
	d := time.Duration(int64(hours * 3600 * 1e9))
	return noon.Add(-d), noon.Add(d)
}
