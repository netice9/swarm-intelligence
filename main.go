package main

import (
	"os"

	"github.com/netice9/swarm-intelligence/api"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "bind",
				Value: ":8080",
			},
		},
		Action: func(c *cli.Context) error {
			return api.Start(c.String("bind"))
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
