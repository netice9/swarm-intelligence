package aggregator

import (
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/netice9/swarm-intelligence/core"
)

var l sync.Mutex
var agentStates = map[string]core.State{}

func NewState(remoteAddr string, s core.State) {
	l.Lock()
	defer l.Unlock()
	agentStates[remoteAddr] = s
}

// remove all agent states that are more than 10 seconds in the past
func init() {
	go func() {
		for ; ; time.Sleep(1 * time.Second) {
			l.Lock()
			stale := []string{}
			n := time.Now()
			for r, s := range agentStates {
				if n.Sub(s.Time) > 90*time.Second {
					stale = append(stale, r)
				}
			}

			for _, r := range stale {
				log.Println("deleting", r)
				delete(agentStates, r)
			}

			l.Unlock()
		}
	}()
}

type volumeList []*types.Volume

func (vl volumeList) hasVolume(name string) bool {
	for _, v := range vl {
		if v.Name == name {
			return true
		}
	}
	return false
}

func State() core.State {
	l.Lock()
	defer l.Unlock()

	s := core.State{
		Time:  time.Now(),
		Stats: map[string]types.Stats{},
	}

	volumes := volumeList{}

	for _, rs := range agentStates {
		s.Tasks = rs.Tasks
		s.Services = rs.Services
		s.Containers = append(s.Containers, rs.Containers...)
		for id, st := range rs.Stats {
			s.Stats[id] = st
		}

		for _, v := range rs.Volumes {
			if !volumes.hasVolume(v.Name) {
				volumes = append(volumes, v)
			}
		}

	}

	s.Volumes = volumes

	return s
}
