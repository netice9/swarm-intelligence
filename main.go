package main

import (
	"fmt"
	"os"

	"github.com/netice9/swarm-intelligence/ui"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	ui.Run(fmt.Sprintf(":%s", port))

}
