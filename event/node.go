package event

import (
	"time"

	"github.com/draganm/emission"
	"github.com/fsouza/go-dockerclient"
)

var Nodes *emission.Emitter = emission.NewEmitter()

var Tasks *emission.Emitter = emission.NewEmitter()

var Services *emission.Emitter = emission.NewEmitter()

func StartCollecting(client *docker.Client) {
	go func() {
		for {
			services, err := client.ListServices(docker.ListServicesOptions{})
			if err != nil {
				panic(err)
			}

			Services.Emit("update", services)

			nodes, err := client.ListNodes(docker.ListNodesOptions{})
			if err != nil {
				panic(err)
			}

			Nodes.Emit("update", nodes)

			tasks, err := client.ListTasks(docker.ListTasksOptions{})
			if err != nil {
				panic(err)
			}

			Tasks.Emit("update", tasks)

			time.Sleep(time.Second)

		}
	}()

}
