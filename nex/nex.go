// Package nex implements a nex station.
//
// For more information about nex, see https://nightfall.city/nex/info/.
package nex

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"text/template"
	"time"
	"blekksprut.net/yofukashi"
)

// A Station serves content from FS.
// Only open at night if Nocturnal is true.
// Uses Latitude to roughly estimate dawn and dusk.
type Station struct {
	FS        fs.FS
	Nocturnal bool
	Latitude  float64
}

// Reads a nex request from rw and tries to serve the matching file.
func (station *Station) Serve(rw io.ReadWriteCloser) error {
	now := time.Now()
	return station.ServeAt(now, rw)
}

// Tries to serve a request at the specific time tm.
// Useful for testing Nocturnal stations.
func (station *Station) ServeAt(tm time.Time, rw io.ReadWriteCloser) error {
	defer rw.Close()

	dawn, dusk := yofukashi.DawnDusk(tm, station.Latitude)

	if station.Nocturnal && tm.Before(dusk) && tm.After(dawn) {
		t, err := template.ParseFS(station.FS, "closed.nex")
		if err != nil {
			d := dusk.Sub(tm)
			var when string
			switch {
			case d.Hours() > 2:
				when = fmt.Sprintf("in about %d hours", int(d.Hours()))
			case d.Hours() > 1:
				when = fmt.Sprintf("in an hour or two")
			case d.Minutes() > 5:
				round := d.Round(5 * time.Minute)
				when = fmt.Sprintf("in about %d minutes", int(round.Minutes()))
			case d.Seconds() > 1:
				when = "soon"
			}
			fmt.Fprintf(rw, "it's still light out. come back %s...", when)
		} else {
			t.Execute(rw, struct{ Dawn, Dusk, Now time.Time }{dawn, dusk, tm})
		}
		return fmt.Errorf("outside opening hours")
	}

	reader := bufio.NewScanner(rw)
	reader.Scan()
	request := reader.Text()

	request = strings.TrimPrefix(request, "/")
	if request == "" || request[len(request)-1] == '/' {
		request = request + "index.nex"
	}

	f, err := station.FS.Open(request)
	if err != nil {
		fmt.Fprintln(rw, "document not found")
		return err
	}
	defer f.Close()

	io.Copy(rw, f)

	return nil
}