package event

import (
	"time"

	"github.com/fsouza/go-dockerclient"
)

var Nodes = NewEmitterAdapter()

var Tasks = NewEmitterAdapter()

var Services = NewEmitterAdapter()

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
