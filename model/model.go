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
	*emission.Emitter
}

var SwarmService *Service

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

	SwarmService = &Service{
		client:    client,
		SwarmInfo: sw,
		Services:  services,
		Emitter:   emission.NewEmitter(),
	}

	go func() {
		for {
			services, err := client.ListServices(docker.ListServicesOptions{})
			if err != nil {
				panic(err)
			}

			if len(services) != len(SwarmService.Services) {
				SwarmService.Services = services
				SwarmService.Emit("services", services)
			}
			time.Sleep(time.Second)

		}
	}()

}
