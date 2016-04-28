package main

import (
	"os"

	"github.com/moznion/floodgate"
)

func main() {
	floodgate.Run(os.Args[1:])
}
