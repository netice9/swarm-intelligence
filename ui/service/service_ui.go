package service

import (
	"encoding/json"
	"fmt"

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

	servceInfo := services.Aggregator.GetServiceInfo(serviceID)

	if servceInfo == nil {
		return reactor.DefaultNotFoundScreenFactory(ctx)
	}

	return &ServiceUI{
		ctx:     ctx,
		ID:      serviceID,
		service: servceInfo,
	}

}

var ui = core.MustParseDisplayModel(`
	<div>
		<div class="page-header">
	  	<h1>Service <span id="serviceName">Name of the service</span>: <small id="serviceID">Subtext for header</small></h1>
		</div>

		<bs.Panel header="Tasks" >
			<bs.ListGroup id="taskList" />
		</bs.Panel>

	  <bs.Panel header="Descriptor" >
			<pre id="text"/>
	  </bs.Panel>
	</div>
`)

var taskItemUI = core.MustParseDisplayModel(`
	<bs.ListGroupItem id="item" header="Heading 1">Some body text</bs.ListGroupItem>
`)

func (s *ServiceUI) Mount() {
	services.Aggregator.OnServiceInfo(s.ID, s.UpdateService)
}

func (s *ServiceUI) render() {
	m := ui.DeepCopy()

	m.SetElementText("serviceName", s.service.Service.Spec.Name)
	m.SetElementText("serviceID", s.service.Service.ID)

	for _, t := range s.service.GetTasks() {
		item := taskItemUI.DeepCopy()
		item.SetElementAttribute("item", "header", t.ID)
		item.SetElementText("item", t.State)
		style := "success"
		if t.State == "failed" {
			style = "danger"
		}
		item.SetElementAttribute("item", "bsStyle", style)
		item.SetElementAttribute("item", "href", fmt.Sprintf("#/task/%s", t.ID))
		m.AppendChild("taskList", item)
	}

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
	services.Aggregator.RemoveServiceInfoListener(s.ID, s.UpdateService)
}
