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
	"time"
)

func serve(rw io.ReadWriteCloser, fs fs.FS) {
	defer rw.Close()

	now := time.Now()
	if now.Hour() >= 7 && now.Hour() < 19 {
		formatted := now.Format("15:04")
		fmt.Fprintf(rw, "it's only %s. come back tonight...", formatted)
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
		go serve(socket, fs)
	}
}
