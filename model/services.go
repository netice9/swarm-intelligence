package model

import (
	"github.com/netice9/swarm-intelligence/event"
	"github.com/netice9/swarm-intelligence/model/services"
)

var Services = services.NewServicesAggregator(event.NewEmitterAdapter())

func init() {
	event.Services.AddListener("update", Services.OnServices)
	event.Tasks.AddListener("update", Services.OnTasks)
}
