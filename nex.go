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

var Now = time.Now

func (nex *Nex) Serve(rw io.ReadWriteCloser) error {
	defer rw.Close()

	now := Now()
	rise, set := sunrise.SunriseSunset(
		nex.Latitude, nex.Longitude,
		now.Year(), now.Month(), now.Day(),
	)

	if nex.Nocturnal && now.Before(set) && now.After(rise) {
		t, err := template.ParseFS(nex.FS, "closed.nex")
		if err != nil {
			d := time.Until(set)
			var when string
			switch {
			case d.Hours() > 1:
				when = fmt.Sprintf("about %d hours", int(d.Hours()))
			case d.Minutes() > 1:
				when = fmt.Sprintf("about %d minutes", int(d.Minutes()))
			case d.Seconds() > 1:
				when = fmt.Sprintf("about %d seconds", int(d.Seconds()))
			default:
				fmt.Fprintf(rw, "we're just about to open, hang on")
				break
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
