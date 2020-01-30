package main

import (
	"edison/internal/pkg/edison"
	"os"
)

func main() {
	meny, err := edison.Collect()
	if err != nil {
		os.Exit(1)
	}

	meny.Render()
}
