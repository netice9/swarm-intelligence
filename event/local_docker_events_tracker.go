package event

import "github.com/fsouza/go-dockerclient"

var DockerEvents = NewEmitterAdapter()

func StartLocalDockerEventsTracking(client *docker.Client) {
	events := make(chan *docker.APIEvents)
	client.AddEventListener(events)
	go func() {
		for evt := range events {
			DockerEvents.Emit("event", evt)
		}
	}()
}
