package node

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/model"
	"github.com/netice9/swarm-intelligence/ui/layout"
	"gitlab.netice9.com/dragan/go-reactor"
	"gitlab.netice9.com/dragan/go-reactor/core"
)

type NodeUI struct {
	ctx  reactor.ScreenContext
	node *swarm.Node
	ID   string
}

func NodeUIFactory(ctx reactor.ScreenContext) reactor.Screen {
	nodeID := ctx.Params["id"]
	node := model.SwarmService.GetNode(nodeID)
	return &NodeUI{
		ctx:  ctx,
		ID:   nodeID,
		node: node,
	}
}

var ui = core.MustParseDisplayModel(`
  <bs.Panel id="mainPanel" >
		<pre id="text"/>
  </bs.Panel>
`)

func (s *NodeUI) Mount() {
	s.render()
}

func (s *NodeUI) render() {
	m := ui.DeepCopy()
	m.SetElementAttribute("mainPanel", "header", fmt.Sprintf("Service %s", s.node.Spec.Name))

	data, err := json.MarshalIndent(s.node, "", "  ")
	if err != nil {
		panic(err)
	}

	m.SetElementText("text", string(data))

	s.ctx.UpdateScreen(&core.DisplayUpdate{Model: layout.WithLayout(m)})
}

func (s *NodeUI) OnUserEvent(evt *core.UserEvent) {
}

func (s *NodeUI) Unmount() {
}
