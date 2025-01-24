// Package nex provides nex station (server) and client functionality.
//
// For more information about nex, see https://nightfall.city/nex/info/.
package nex

import (
	"blekksprut.net/yofukashi"
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/url"
	"strings"
	"text/template"
	"time"
)

// A Station serves content from FS.
// Only open at night if Nocturnal is true.
// Uses Latitude to roughly estimate dawn and dusk.
type Station struct {
	FS        fs.FS
	Nocturnal bool
	Latitude  float64
}

// A Response represents a response from a nex station.
// It implements the io.Reader interface.
type Response struct {
	Raw  io.Reader
	Conn net.Conn
}

// Read reads up to len(b) bytes into b.
func (r *Response) Read(b []byte) (int, error) {
	return r.Raw.Read(b)
}

// Close the connection
func (r *Response) Close() {
	r.Conn.Close()
}

// Request makes a nex request to rawURL.
func Request(ctx context.Context, rawURL string) (*Response, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if url.Port() == "" {
		url.Host = url.Host + ":1900"
	}
	timeout, _ := time.ParseDuration("30s")
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", url.Host)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, url.Path)
	reader := bufio.NewReader(conn)
	return &Response{Raw: reader, Conn: conn}, nil
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
