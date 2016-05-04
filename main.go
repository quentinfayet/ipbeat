package main

import (
	"ipbeat/beater"
	"os"

	"github.com/elastic/beats/libbeat/beat"
)

// Name defines beat's name
var Name = "ipbeat"

func main() {
	err := beat.Run(Name, "", beater.New())

	if err != nil {
		os.Exit(1)
	}

}
