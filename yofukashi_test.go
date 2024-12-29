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

func init() {
	yofukashi.Now = midnight
}

func TestNex(t *testing.T) {
	nex := yofukashi.Nex{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/README.gmi"), Writer: io.Discard}
	err := nex.Serve(req)
	if err != nil {
		t.Errorf("should succeed")
	}
}

func TestIndex(t *testing.T) {
	nex := yofukashi.Nex{FS: os.DirFS("nex")}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := nex.Serve(req)
	if err != nil {
		t.Errorf("should serve up the index")
	}
}

func TestMissingIndex(t *testing.T) {
	nex := yofukashi.Nex{FS: os.DirFS(".")}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := nex.Serve(req)
	if err == nil {
		t.Errorf("no index.nex, should fail")
	}
}

func TestHours(t *testing.T) {
	yofukashi.Now = midday
	nex := yofukashi.Nex{os.DirFS("."), true, 35.6764, 139.6500}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := nex.Serve(req)
	if err == nil {
		t.Errorf("outside opening hours, should fail")
	}
}

func TestClosingTemplate(t *testing.T) {
	yofukashi.Now = midday
	nex := yofukashi.Nex{os.DirFS("nex"), true, 35.6764, 139.6500}
	req := request{Reader: strings.NewReader("/"), Writer: io.Discard}
	err := nex.Serve(req)
	if err == nil {
		t.Errorf("outside opening hours, should fail")
	}
}
