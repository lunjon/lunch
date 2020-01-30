package main

import (
	"github.com/lunjon/lunch/internal/pkg/edison"
	"os"
)

func main() {
	meny, err := edison.Collect()
	if err != nil {
		os.Exit(1)
	}

	meny.Render()
}
