package main

import (
	"blekksprut.net/yofukashi"
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"strings"
	"text/template"
	"time"
)

func serve(rw io.ReadWriteCloser, fs fs.FS, hours bool) {
	defer rw.Close()

	now := time.Now()
	if hours && now.Hour() >= 7 && now.Hour() < 19 {
		t, err := template.ParseFS(fs, "closed.nex")
		if err != nil {
			formatted := now.Format("15:04")
			fmt.Fprintf(rw, "it's only %s. come back tonight...", formatted)
		} else {
			t.Execute(rw, now)
		}
		return
	}

	reader := bufio.NewScanner(rw)
	reader.Scan()
	request := reader.Text()

	request = strings.TrimPrefix(request, "/")
	if request == "" || request[len(request)-1] == '/' {
		request = request + "index.nex"
	}

	f, err := fs.Open(request)
	if err != nil {
		fmt.Fprintln(rw, "document not found")
		return
	}
	defer f.Close()

	io.Copy(rw, f)
}

func main() {
	r := flag.String("r", "/var/nex", "root directory")
	v := flag.Bool("v", false, "version")
	a := flag.Bool("a", false, "keep open around the clock")
	flag.Parse()

	if *v {
		fmt.Println(os.Args[0], yofukashi.Version)
		os.Exit(0)
	}

	fs := os.DirFS(*r)

	Lockdown(*r)

	server, err := net.Listen("tcp", ":1900")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		socket, err := server.Accept()
		if err != nil {
			log.Println(err)
		}
		go serve(socket, fs, !*a)
	}
}
