package main

import (
	"os"

	"github.com/scottbrown/scruffy"
)

func main() {
	if err := scruffy.Execute(); err != nil {
		os.Exit(1)
	}
}