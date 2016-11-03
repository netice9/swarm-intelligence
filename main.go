package main

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
	"github.com/netice9/swarm-intelligence/event"
	"github.com/netice9/swarm-intelligence/ui"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	event.StartCollecting(client)

	ui.Run(fmt.Sprintf(":%s", port))

}
