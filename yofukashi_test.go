package yofukashi_test

import (
	"blekksprut.net/yofukashi"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

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
	station := yofukashi.Station{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/README.gmi"), Writer: io.Discard}
	station.Serve(req)
}

func TestStation(t *testing.T) {
	station := yofukashi.Station{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/README.gmi"), Writer: io.Discard}
	err := station.ServeAt(midnight(), req)
	if err != nil {
		t.Errorf("should succeed")
	}
}

func TestIndex(t *testing.T) {
	station := yofukashi.Station{FS: os.DirFS("nex")}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midnight(), req)
	if err != nil {
		t.Errorf("should serve up the index")
	}
}

func TestMissingIndex(t *testing.T) {
	station := yofukashi.Station{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midnight(), req)
	if err == nil {
		t.Errorf("no index.nex, should fail")
	}
}

func TestHours(t *testing.T) {
	station := yofukashi.Station{os.DirFS("."), true, 35.6764}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midday(), req)
	if err == nil {
		t.Errorf("outside opening hours, should fail")
	}
}

func TestClosingTemplate(t *testing.T) {
	station := yofukashi.Station{os.DirFS("nex"), true, 35.6764}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := station.ServeAt(midday(), req)
	if err == nil {
		t.Errorf("outside opening hours, should fail")
	}
}

func TestOpeningEstimates(t *testing.T) {
	station := yofukashi.Station{os.DirFS("."), true, 35.6764}
	var res strings.Builder
	req := request{Reader: strings.NewReader("/"), Writer: &res}
	now := time.Now()
	_, dusk := yofukashi.DawnDusk(now, 35.6764)

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
