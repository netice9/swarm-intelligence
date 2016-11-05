package stats

import (
	"log"

	"github.com/fsouza/go-dockerclient"
	"github.com/netice9/swarm-intelligence/event"
)

func StartTracking(client *docker.Client) {

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

	eventsChannel := make(chan *docker.APIEvents)

	err = client.AddEventListener(eventsChannel)
	if err != nil {
		panic(err)
	}

	go func() {
		for evt := range eventsChannel {
			if evt.Status == "start" {
				go trackContainer(evt.ID, client)
			}
		}
	}()

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
		event.ContainerStats.Emit("stats", update)
	}

}
