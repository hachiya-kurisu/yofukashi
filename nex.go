package yofukashi

import (
	"bufio"
	"fmt"
	"github.com/nathan-osman/go-sunrise"
	"io"
	"io/fs"
	"strings"
	"text/template"
	"time"
)

type Nex struct {
	FS        fs.FS
	Nocturnal bool
	Latitude  float64
	Longitude float64
}

func (nex *Nex) Serve(rw io.ReadWriteCloser) error {
	now := time.Now()
	return nex.ServeAt(now, rw)
}

func (nex *Nex) SunriseSunset(tm time.Time) (time.Time, time.Time) {
	return sunrise.SunriseSunset(
		nex.Latitude, nex.Longitude,
		tm.Year(), tm.Month(), tm.Day(),
	)
}

func (nex *Nex) ServeAt(tm time.Time, rw io.ReadWriteCloser) error {
	defer rw.Close()

	rise, set := nex.SunriseSunset(tm)

	if nex.Nocturnal && tm.Before(set) && tm.After(rise) {
		t, err := template.ParseFS(nex.FS, "closed.nex")
		if err != nil {
			d := set.Sub(tm)
			var when string
			switch {
			case d.Hours() > 1:
				when = fmt.Sprintf("about %d hours", int(d.Hours()))
			case d.Minutes() > 1:
				when = fmt.Sprintf("about %d minutes", int(d.Minutes()))
			case d.Seconds() > 1:
				when = "a minute"
			}
			fmt.Fprintf(rw, "it's still light out. come back in %s...", when)
		} else {
			t.Execute(rw, set)
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

	f, err := nex.FS.Open(request)
	if err != nil {
		fmt.Fprintln(rw, "document not found")
		return err
	}
	defer f.Close()

	io.Copy(rw, f)

	return nil
}
