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
	// currentState := State{
	// 	Time: time.Now(),
	// }

	go func() {
		for {
			newState := State{
				Time: time.Now(),
			}
			sl, err := c.ServiceList(context.Background(), types.ServiceListOptions{})
			if err != nil {
				log.Printf("Error fetching services: %s", err.Error())
			}
			newState.Services = sl
			currentState.Store(newState)
			time.Sleep(time.Second)
		}
	}()
}
