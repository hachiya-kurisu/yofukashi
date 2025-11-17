// higure is a nocturnal nex server.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"blekksprut.net/yofukashi"
	"blekksprut.net/yofukashi/nex"
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

	root, err := os.OpenRoot(*r)
	if err != nil {
		log.Fatal(err)
	}

	Lockdown(*r)

	server, err := net.Listen("tcp", ":1900")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	station := nex.Station{FS: root.FS(), Nocturnal: !*a, Latitude: *lat}
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
		go station.Serve(socket)
	}
}
