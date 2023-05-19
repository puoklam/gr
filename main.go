package main

import (
	"github.com/puoklam/gr/cmd"
	"github.com/puoklam/gr/log"
)

func main() {
	defer log.Clear()
	if err := cmd.Exec(); err != nil {
		log.Fatal(err)
	}
}
