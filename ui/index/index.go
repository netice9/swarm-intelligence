package index

import (
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/model"
	"github.com/netice9/swarm-intelligence/model/services"
	"github.com/netice9/swarm-intelligence/ui/layout"
	"gitlab.netice9.com/dragan/go-reactor"
	"gitlab.netice9.com/dragan/go-reactor/core"
)

type Index struct {
	ctx      reactor.ScreenContext
	services services.ServiceList
	nodes    []swarm.Node
	tasks    []swarm.Task
}

func IndexFactory(ctx reactor.ScreenContext) reactor.Screen {
	return &Index{
		ctx: ctx,
	}
}

var ui = core.MustParseDisplayModel(`
	<div>
	  <bs.Panel header="Services">
	    <bs.ListGroup id="services"/>
	  </bs.Panel>
		<bs.Panel header="Nodes">
	    <bs.ListGroup id="nodes"/>
	  </bs.Panel>
		<bs.Panel header="Tasks">
	    <bs.ListGroup id="tasks"/>
	  </bs.Panel>
	</div>
`)

var serviceListItemUI = core.MustParseDisplayModel(`
  <bs.ListGroupItem id="service"/>
`)

var nodeListItemUI = core.MustParseDisplayModel(`
  <bs.ListGroupItem id="node"/>
`)

var taskListItemUI = core.MustParseDisplayModel(`
  <bs.ListGroupItem id="task"/>
`)

func (i *Index) Mount() {
	i.services = model.Services.ServiceList()
	i.nodes = model.SwarmService.Nodes
	i.tasks = model.SwarmService.Tasks
	i.render()
	model.SwarmService.AddListener("services", i.OnServices)
	model.SwarmService.AddListener("nodes", i.OnNodes)
	model.SwarmService.AddListener("tasks", i.OnTasks)
}

func (i *Index) render() {
	m := ui.DeepCopy()
	// m.SetElementText(id string, text string)
	for _, s := range i.services {
		name := s.Name
		item := serviceListItemUI.DeepCopy()
		item.SetElementText("service", name)
		item.SetElementAttribute("service", "href", fmt.Sprintf("#/service/%s", s.ID))
		m.AppendChild("services", item)
	}

	for _, n := range i.nodes {

		item := nodeListItemUI.DeepCopy()
		item.SetElementText("node", n.ID)
		item.SetElementAttribute("node", "href", fmt.Sprintf("#/node/%s", n.ID))
		m.AppendChild("nodes", item)
	}

	for _, t := range i.tasks {

		item := taskListItemUI.DeepCopy()
		item.SetElementText("task", t.ID)
		item.SetElementAttribute("task", "href", fmt.Sprintf("#/task/%s", t.ID))
		m.AppendChild("tasks", item)
	}

	i.ctx.UpdateScreen(&core.DisplayUpdate{Model: layout.WithLayout(m)})
}

func (i *Index) OnServices(services services.ServiceList) {
	i.services = services
	i.render()
}

func (i *Index) OnNodes(nodes []swarm.Node) {
	i.nodes = nodes
	i.render()
}

func (i *Index) OnTasks(tasks []swarm.Task) {
	i.tasks = tasks
	i.render()
}

func (i *Index) OnUserEvent(evt *core.UserEvent) {
}

func (i *Index) Unmount() {
	model.SwarmService.RemoveListener("services", i.OnServices)
	model.SwarmService.RemoveListener("nodes", i.OnNodes)
}
