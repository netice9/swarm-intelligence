package ui

import (
	"github.com/netice9/swarm-intelligence/ui/index"
	"gitlab.netice9.com/dragan/go-reactor"
)

func Run(bnd string) {
	r := reactor.New()
	r.AddScreen("/", index.IndexFactory)
	r.Serve(bnd)
}
