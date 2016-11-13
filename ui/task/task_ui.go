package task

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/model"
	"github.com/netice9/swarm-intelligence/ui/layout"
	"gitlab.netice9.com/dragan/go-reactor"
	"gitlab.netice9.com/dragan/go-reactor/core"
)

type TaskUI struct {
	ctx     reactor.ScreenContext
	task    *swarm.Task
	ID      string
	running bool
}

func TaskUIFactory(ctx reactor.ScreenContext) reactor.Screen {
	taskID := ctx.Params["id"]
	task := model.SwarmService.GetTask(taskID)

	if task == nil {
		return reactor.DefaultNotFoundScreenFactory(ctx)
	}

	return &TaskUI{
		ctx:  ctx,
		ID:   taskID,
		task: task,
	}

}

var ui = core.MustParseDisplayModel(`
  <bs.Panel id="mainPanel" >
		<pre id="text"/>
  </bs.Panel>
`)

func (t *TaskUI) Mount() {
	t.running = true
	model.SwarmService.AddListener(fmt.Sprintf("task/%s", t.ID), t.OnTask)
	t.render()
}

func (t *TaskUI) OnTask(ts swarm.Task) {
	// fmt.Println("on task", ts)
	t.task = &ts
	t.render()
}

func (t *TaskUI) render() {
	if !t.running {
		return
	}
	m := ui.DeepCopy()
	m.SetElementAttribute("mainPanel", "header", fmt.Sprintf("Task %s", t.task.Name))

	data, err := json.MarshalIndent(t.task, "", "  ")
	if err != nil {
		panic(err)
	}

	m.SetElementText("text", string(data))

	t.ctx.UpdateScreen(&core.DisplayUpdate{Model: layout.WithLayout(m)})
}

func (t *TaskUI) OnUserEvent(evt *core.UserEvent) {
}

func (t *TaskUI) Unmount() {
	t.running = false
	model.SwarmService.RemoveListener(fmt.Sprintf("task/%s", t.ID), t.OnTask)
}
