package service

import (
	"encoding/json"
	"fmt"

	"github.com/netice9/swarm-intelligence/model"
	"github.com/netice9/swarm-intelligence/model/services"
	"github.com/netice9/swarm-intelligence/ui/layout"
	"gitlab.netice9.com/dragan/go-reactor"
	"gitlab.netice9.com/dragan/go-reactor/core"
)

type ServiceUI struct {
	ctx     reactor.ScreenContext
	service *services.ServiceInfo
	ID      string
}

func ServiceUIFactory(ctx reactor.ScreenContext) reactor.Screen {
	serviceID := ctx.Params["id"]
	return &ServiceUI{
		ctx:     ctx,
		ID:      serviceID,
		service: model.Services.GetServiceInfo(serviceID),
	}

}

var ui = core.MustParseDisplayModel(`
  <bs.Panel id="mainPanel" >
		<pre id="text"/>
  </bs.Panel>
`)

func (s *ServiceUI) Mount() {
	s.render()
	model.Services.AddListener(fmt.Sprintf("update/%s", s.ID), s.UpdateService)
}

func (s *ServiceUI) render() {
	m := ui.DeepCopy()
	m.SetElementAttribute("mainPanel", "header", fmt.Sprintf("Service %s", s.service.Service.Spec.Name))

	data, err := json.MarshalIndent(s.service, "", "  ")
	if err != nil {
		panic(err)
	}

	m.SetElementText("text", string(data))

	s.ctx.UpdateScreen(&core.DisplayUpdate{Model: layout.WithLayout(m)})
}

func (s *ServiceUI) OnUserEvent(evt *core.UserEvent) {
}

func (s *ServiceUI) UpdateService(info *services.ServiceInfo) {
	s.service = info
	s.render()
}

func (s *ServiceUI) Unmount() {
	model.Services.RemoveListener(fmt.Sprintf("update/%s", s.ID), s.UpdateService)
}
