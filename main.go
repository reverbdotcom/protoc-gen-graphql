package main

import (
	"bytes"
	"log"

	"github.com/pseudomuto/protokit"
)

func main() {
	if err := protokit.RunPlugin(&plugin{out: &bytes.Buffer{}}); err != nil {
		log.Fatal(err)
	}
}
