// hoshikuzu is a client for the nex protocol
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"blekksprut.net/yofukashi"
	"blekksprut.net/yofukashi/nex"
	"github.com/blacktop/go-termimg"
)

func displayImage(path string) {
	img, err := termimg.Open(path)
	if err != nil {
		panic(err)
	}
	is, err := img.Render()
	if err != nil {
		panic(err)
	}
	fmt.Println(is)
}

func main() {
	v := flag.Bool("v", false, "version")

	flag.Parse()

	if *v {
		fmt.Println(os.Args[0], yofukashi.Version)
		os.Exit(0)
	}

	for _, arg := range flag.Args() {
		if !strings.HasPrefix(arg, "nex://") {
			arg = "nex://" + arg
		}
		r, err := nex.Request(context.Background(), arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		defer r.Close()
		base := filepath.Base(arg)
		switch strings.ToLower(filepath.Ext(base)) {
		case ".jpg", ".jpeg", ".png":
			f, err := os.CreateTemp("", "*"+base)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				io.Copy(f, r)
				path := f.Name()
				protocol := termimg.DetectProtocol()
				switch protocol {
				case termimg.ITerm2, termimg.Kitty:
					displayImage(path)
				default:
					exec.Command("open", path).Run()
				}
			}
		default:
			io.Copy(os.Stdout, r)
		}
	}
}
