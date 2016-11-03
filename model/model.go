package model

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/draganm/emission"
	"github.com/fsouza/go-dockerclient"
)

type Service struct {
	client    *docker.Client
	SwarmInfo swarm.Info
	Services  []swarm.Service
	Nodes     []swarm.Node
	Tasks     []swarm.Task
	*emission.Emitter
}

var SwarmService *Service

func (s *Service) GetService(serviceID string) *swarm.Service {
	for _, srv := range s.Services {
		if srv.ID == serviceID {
			return &srv
		}
	}
	return nil
}

func (s *Service) GetTask(taskID string) *swarm.Task {
	for _, tsk := range s.Tasks {
		if tsk.ID == taskID {
			return &tsk
		}
	}
	return nil
}

func (s *Service) GetNode(nodeID string) *swarm.Node {
	for _, nod := range s.Nodes {
		if nod.ID == nodeID {
			return &nod
		}
	}
	return nil
}

func init() {

	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	info, err := client.Info()
	sw := info.Swarm

	services, err := client.ListServices(docker.ListServicesOptions{})
	if err != nil {
		panic(err)
	}

	nodes, err := client.ListNodes(docker.ListNodesOptions{})
	if err != nil {
		panic(err)
	}

	tasks, err := client.ListTasks(docker.ListTasksOptions{})
	if err != nil {
		panic(err)
	}

	SwarmService = &Service{
		client:    client,
		SwarmInfo: sw,
		Services:  services,
		Emitter:   emission.NewEmitter(),
		Nodes:     nodes,
		Tasks:     tasks,
	}

	go func() {
		for {
			services, err := client.ListServices(docker.ListServicesOptions{})
			if err != nil {
				panic(err)
			}

			SwarmService.Services = services

			if len(services) != len(SwarmService.Services) {
				SwarmService.Emit("services", services)
			}

			nodes, err := client.ListNodes(docker.ListNodesOptions{})
			if err != nil {
				panic(err)
			}

			SwarmService.Nodes = nodes

			if len(nodes) != len(SwarmService.Nodes) {
				SwarmService.Emit("nodes", nodes)
			}

			tasks, err := client.ListTasks(docker.ListTasksOptions{})
			if err != nil {
				panic(err)
			}

			SwarmService.Tasks = tasks

			if len(tasks) != len(SwarmService.Tasks) {
				SwarmService.Emit("tasks", tasks)
			}

			for _, t := range tasks {
				SwarmService.Emit(fmt.Sprintf("task/%s", t.ID), t)
			}

			time.Sleep(time.Second)

		}
	}()

}
