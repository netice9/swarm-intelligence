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
	</div>
`)

var serviceListItemUI = core.MustParseDisplayModel(`
  <bs.ListGroupItem id="service">
	<div>
		CPU: <span id="cpu"/>
		Mem: <span id="memory"/>
	</div>
	</bs.ListGroupItem>
`)

var nodeListItemUI = core.MustParseDisplayModel(`
  <bs.ListGroupItem id="node"/>
`)

func (i *Index) Mount() {
	i.nodes = model.SwarmService.Nodes
	services.Aggregator.OnServiceList(i.onServiceList)
	model.SwarmService.AddListener("nodes", i.OnNodes)
}

func (i *Index) render() {
	m := ui.DeepCopy()
	for _, s := range i.services {
		name := s.Name
		item := serviceListItemUI.DeepCopy()
		// item.SetElementText("service", fmt.Sprintf("%s Mem: %d, CPU: %d", name, s.Memory, s.CPU))
		item.SetElementAttribute("service", "header", name)
		item.SetElementText("cpu", fmt.Sprintf("%.1f%%", float64(s.CPU)/1e7))
		item.SetElementText("memory", fmt.Sprintf("%.1f MB", float64(s.CPU)/(1024*1024)))
		item.SetElementAttribute("service", "href", fmt.Sprintf("#/service/%s", s.ID))
		m.AppendChild("services", item)
	}

	for _, n := range i.nodes {

		item := nodeListItemUI.DeepCopy()
		item.SetElementText("node", n.ID)
		item.SetElementAttribute("node", "href", fmt.Sprintf("#/node/%s", n.ID))
		m.AppendChild("nodes", item)
	}

	i.ctx.UpdateScreen(&core.DisplayUpdate{Model: layout.WithLayout(m)})
}

func (i *Index) onServiceList(services services.ServiceList) {
	i.services = services
	i.render()
}

func (i *Index) OnNodes(nodes []swarm.Node) {
	i.nodes = nodes
	i.render()
}

func (i *Index) OnUserEvent(evt *core.UserEvent) {
}

func (i *Index) Unmount() {
	services.Aggregator.RemoveServiceListListener(i.onServiceList)
	model.SwarmService.RemoveListener("nodes", i.OnNodes)
}
