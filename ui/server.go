package ui

import (
	"github.com/netice9/swarm-intelligence/ui/index"
	"github.com/netice9/swarm-intelligence/ui/node"
	"github.com/netice9/swarm-intelligence/ui/service"
	"github.com/netice9/swarm-intelligence/ui/task"
	"gitlab.netice9.com/dragan/go-reactor"
)

func Run(bnd string) {
	r := reactor.New()
	r.AddScreen("/", index.IndexFactory)
	r.AddScreen("/service/:id", service.ServiceUIFactory)
	r.AddScreen("/node/:id", node.NodeUIFactory)
	r.AddScreen("/task/:id", task.TaskUIFactory)
	r.Serve(bnd)
}
