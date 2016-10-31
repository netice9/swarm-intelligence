package index

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/model"
	"github.com/netice9/swarm-intelligence/ui/layout"
	"gitlab.netice9.com/dragan/go-reactor"
	"gitlab.netice9.com/dragan/go-reactor/core"
)

type Index struct {
	ctx      reactor.ScreenContext
	services []swarm.Service
}

func IndexFactory(ctx reactor.ScreenContext) reactor.Screen {
	return &Index{
		ctx: ctx,
	}
}

var ui = core.MustParseDisplayModel(`
  <bs.Panel header="Services">
    <bs.ListGroup id="services"/>
  </bs.Panel>
`)

var serviceListItemUI = core.MustParseDisplayModel(`
  <bs.ListGroupItem id="service"/>
`)

var boardLinkUI = core.MustParseDisplayModel(`<bs.ListGroupItem id="link"/>`)

func (i *Index) Mount() {
	i.services = model.SwarmService.Services
	i.render()
	model.SwarmService.AddListener("services", i.OnServices)
}

func (i *Index) render() {
	m := ui.DeepCopy()
	// m.SetElementText(id string, text string)
	for _, s := range i.services {
		name := s.Spec.Name
		item := serviceListItemUI.DeepCopy()
		item.SetElementText("service", name)
		m.AppendChild("services", item)
	}
	i.ctx.UpdateScreen(&core.DisplayUpdate{Model: layout.WithLayout(m)})
}

func (i *Index) OnServices(services []swarm.Service) {

	i.services = services
	i.render()
}

func (i *Index) OnUserEvent(evt *core.UserEvent) {
}

func (i *Index) Unmount() {
	model.SwarmService.RemoveListener("services", i.OnServices)
}
