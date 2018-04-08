package main

import (
	"os"

	"github.com/netice9/swarm-intelligence/agent"
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
			&cli.StringFlag{
				Name:  "agent-bind",
				Value: ":9000",
			},
		},
		Action: func(c *cli.Context) error {
			return api.Start(
				c.String("bind"),
				c.String("agent-bind"),
			)
		},
		Commands: []*cli.Command{
			&cli.Command{

				Name: "agent",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "remote",
						Value: "http://head:9000",
					},
				},
				Action: func(c *cli.Context) error {
					return agent.Run(c.String("remote"))
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
