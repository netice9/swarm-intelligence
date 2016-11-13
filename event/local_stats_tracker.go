package event

import (
	"log"

	"github.com/fsouza/go-dockerclient"
)

type StatsForContainer struct {
	Stats       *docker.Stats
	ContainerID string
}

var ContainerStats = NewEmitterAdapter()

func StartTrackingLocalContainerStats(client *docker.Client) {

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		if c.State == "running" {
			go trackContainer(c.ID, client)
		}
	}

	DockerEvents.On("event", func(evt *docker.APIEvents) {
		if evt.Status == "start" {
			go trackContainer(evt.ID, client)
		}
	})

}

func trackContainer(id string, client *docker.Client) {
	statsChan := make(chan *docker.Stats)

	go func() {
		err := client.Stats(docker.StatsOptions{Stream: true, ID: id, Stats: statsChan})
		if err != nil {
			log.Printf("Error getting status of container %s: %s", id, err)
			return
		}
	}()

	for update := range statsChan {
		ContainerStats.Emit("stats", StatsForContainer{Stats: update, ContainerID: id})
	}

}
