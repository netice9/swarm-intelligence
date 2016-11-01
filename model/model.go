package model

import (
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/docker/docker/api/types/swarm"
	"github.com/fsouza/go-dockerclient"
)

type Service struct {
	client    *docker.Client
	SwarmInfo swarm.Info
	Services  []swarm.Service
	Nodes     []swarm.Node
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

	SwarmService = &Service{
		client:    client,
		SwarmInfo: sw,
		Services:  services,
		Emitter:   emission.NewEmitter(),
		Nodes:     nodes,
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

			time.Sleep(time.Second)

		}
	}()

}
