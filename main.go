package main

import (
	"os"

	"github.com/mhrynenko/jwt_service/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
