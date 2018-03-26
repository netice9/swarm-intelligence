package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type State struct {
	Time       time.Time              `json:"time"`
	Services   []swarm.Service        `json:"services"`
	Tasks      []swarm.Task           `json:"tasks"`
	Containers []types.Container      `json:"containers"`
	Stats      map[string]types.Stats `json:"stats"`
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

			cl, err := c.ContainerList(context.Background(), types.ContainerListOptions{})
			if err != nil {
				log.Printf("Error fetching containers: %s", err.Error())
			}

			stats := map[string]types.Stats{}

			for _, con := range cl {
				cs, err := c.ContainerStats(context.Background(), con.ID, false)
				if err != nil {
					log.Printf("Error fetching containers: %s", err.Error())
				}

				st := types.Stats{}
				err = json.NewDecoder(cs.Body).Decode(&st)
				if err != nil {
					log.Printf("Error fetching containers: %s", err.Error())
				}
				stats[con.ID] = st
				cs.Body.Close()

			}

			// c.ContainerStats(ctx, containerID, stream)

			newState := State{
				Time:       time.Now(),
				Services:   sl,
				Tasks:      tl,
				Containers: cl,
				Stats:      stats,
			}
			currentState.Store(newState)
			time.Sleep(time.Second)
		}
	}()
}
