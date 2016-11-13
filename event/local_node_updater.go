package event

import "github.com/fsouza/go-dockerclient"

var Nodes = NewEmitterAdapter()

var Tasks = NewEmitterAdapter()

var Services = NewEmitterAdapter()

func StartCollecting(client *docker.Client) {
	Time.On("1sec", func() {
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
	})
}
