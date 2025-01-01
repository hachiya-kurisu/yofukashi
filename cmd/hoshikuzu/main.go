// hoshikuzu is a client for the nex protocol
package main

import (
	"blekksprut.net/yofukashi"
	"blekksprut.net/yofukashi/nex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	v := flag.Bool("v", false, "version")

	flag.Parse()

	if *v {
		fmt.Println(os.Args[0], yofukashi.Version)
		os.Exit(0)
	}

	for _, arg := range flag.Args() {
		r, err := nex.Request(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			base := filepath.Base(arg)
			switch strings.ToLower(filepath.Ext(base)) {
			case ".jpg", ".jpeg", ".png", ".gif":
				f, err := os.CreateTemp("", "*"+base)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					io.Copy(f, r)
					exec.Command("open", f.Name()).Run()
				}
			default:
				io.Copy(os.Stdout, r)
			}
		}
	}
}
