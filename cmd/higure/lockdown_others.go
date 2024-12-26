//go:build !openbsd

package main

import (
	"log"
)

func Lockdown(path string) {
	log.Println("no pledge/unveil 😭")
	return
}
