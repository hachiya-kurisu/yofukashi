package main

import (
	"blekksprut.net/yofukashi"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
	"os"
)

func main() {
	r := flag.String("r", "/var/nex", "root directory")
	v := flag.Bool("v", false, "version")
	a := flag.Bool("a", false, "keep open around the clock")
	lat := flag.Float64("lat", 35.68, "latitude")

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

	nex := yofukashi.Nex{FS: fs, Nocturnal: !*a, Latitude: *lat}
	log.Printf("listening on :1900")
	if !*a {
		now := time.Now()
		dawn, dusk := yofukashi.DawnDusk(now, *lat)
		log.Printf("%s to %s", dusk.Format("15:04"), dawn.Format("15:04"))
	}
	for {
		socket, err := server.Accept()
		if err != nil {
			log.Println(err)
		}
		go nex.Serve(socket)
	}
}
