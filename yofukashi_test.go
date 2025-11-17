package yofukashi_test

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"blekksprut.net/yofukashi"
	"blekksprut.net/yofukashi/nex"
)

const lat = 35.6764

type request struct {
	io.Writer
	io.Reader
}

func (request) Close() error {
	return nil
}

func midnight() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func midday() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 12, 1, 0, 0, t.Location())
}

func TestServe(t *testing.T) {
	station := nex.Station{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/README.gmi"), Writer: io.Discard}
	station.Serve(req)
}

func TestStation(t *testing.T) {
	station := nex.Station{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/README.gmi"), Writer: io.Discard}
	err := station.ServeAt(midnight(), req)
	if err != nil {
		t.Errorf("should succeed")
	}
}

func TestIndex(t *testing.T) {
	station := nex.Station{FS: os.DirFS("station")}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midnight(), req)
	if err != nil {
		t.Errorf("should serve up the index")
	}
}

func TestMissingIndex(t *testing.T) {
	station := nex.Station{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midnight(), req)
	if err == nil {
		t.Errorf("no index.nex, should fail")
	}
}

func TestHours(t *testing.T) {
	station := nex.Station{os.DirFS("."), true, lat}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midday(), req)
	if err == nil {
		t.Errorf("outside opening hours, should fail")
	}
}

func TestClosingTemplate(t *testing.T) {
	station := nex.Station{os.DirFS("station"), true, lat}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midday(), req)
	if err == nil {
		t.Errorf("outside opening hours, should fail")
	}
}

func TestDaytime(t *testing.T) {
	midday, _ := time.Parse("15:04", "12:00")
	if yofukashi.Nighttime(midday, lat) {
		t.Errorf("12:00 should be considered daytime")
	}
}

func TestNighttime(t *testing.T) {
	evening, _ := time.Parse("15:04", "21:00")
	if yofukashi.Daytime(evening, lat) {
		t.Errorf("21:00 should be considered nighttime")
	}
}

func TestOpeningEstimates(t *testing.T) {
	station := nex.Station{os.DirFS("."), true, lat}
	var res strings.Builder
	req := request{Reader: strings.NewReader("/"), Writer: &res}
	now := time.Now()
	_, dusk := yofukashi.DawnDusk(now, lat)

	t.Run("Hours", func(t *testing.T) {
		d, _ := time.ParseDuration("-5h")
		station.ServeAt(dusk.Add(d), req)
		if !strings.Contains(res.String(), "5 hours") {
			t.Errorf("failed to estimate number of hours until opening")
		}
	})
	t.Run("AFewHours", func(t *testing.T) {
		d, _ := time.ParseDuration("-90m")
		station.ServeAt(dusk.Add(d), req)
		if !strings.Contains(res.String(), "an hour or two") {
			t.Errorf("failed to estimate number of hours until opening")
		}
	})
	t.Run("Minutes", func(t *testing.T) {
		d, _ := time.ParseDuration("-11m")
		station.ServeAt(dusk.Add(d), req)
		if !strings.Contains(res.String(), "10 minutes") {
			t.Errorf("failed to estimate number of minutes until opening")
		}
	})
	t.Run("Soon", func(t *testing.T) {
		d, _ := time.ParseDuration("-3s")
		station.ServeAt(dusk.Add(d), req)
		if !strings.Contains(res.String(), "soon") {
			t.Errorf("failed to estimate number of seconds until opening")
		}
	})
}

func TestRequest(t *testing.T) {
	u := os.Getenv("YOFUKASHI_TEST_URL")
	if u != "" {
		r, err := nex.Request(context.Background(), u)
		if err != nil {
			t.Errorf("failed to get test url")
		} else {
			defer r.Close()
			ioutil.ReadAll(r)
		}

		_, err = nex.Request(context.Background(), strings.Replace(u, "nex", "", 1))
		if err == nil {
			t.Errorf("request to invalid url should fail")
		}

		_, err = nex.Request(context.Background(), "nex://broken:123456/")
		if err == nil {
			t.Errorf("request to invalid url should fail")
		}
	}
}
