package core

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type State struct {
	Time     time.Time       `json:"time"`
	Services []swarm.Service `json:"services"`
	Tasks    []swarm.Task    `json:"tasks"`
}

var currentState atomic.Value

func CurrentState() State {
	return currentState.Load().(State)
}

func init() {

	currentState.Store(State{Time: time.Now()})

	c, err := client.NewEnvClient()
	if err != nil {
		panic(fmt.Errorf("Could not intialize docker client: %s", err.Error()))
	}

	go func() {
		for {
			sl, err := c.ServiceList(context.Background(), types.ServiceListOptions{})
			if err != nil {
				log.Printf("Error fetching services: %s", err.Error())
			}

			tl, err := c.TaskList(context.Background(), types.TaskListOptions{})

			if err != nil {
				log.Printf("Error fetching tasks: %s", err.Error())
			}

			newState := State{
				Time:     time.Now(),
				Services: sl,
				Tasks:    tl,
			}
			currentState.Store(newState)
			time.Sleep(time.Second)
		}
	}()
}
