package services

import (
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/model"
)

type ServiceStatus struct {
	Name string
	ID   string
}

type ServiceList []ServiceStatus

func (sl ServiceList) Len() int {
	return len(sl)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (sl ServiceList) Less(i, j int) bool {
	if sl[i].Name == sl[j].Name {
		return sl[i].ID < sl[j].ID
	}
	return sl[i].Name < sl[j].Name
}

// Swap swaps the elements with indexes i and j.
func (sl ServiceList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

type ServicesAggregator struct {
	sync.Mutex
	// current []swarm.Service
	current map[string]swarm.Service
	model.EventEmitter
	serviceList ServiceList
}

func (sa *ServicesAggregator) ServiceList() []ServiceStatus {
	sa.Lock()
	defer sa.Unlock()
	return sa.serviceList
}

func NewServicesAggregator(ee model.EventEmitter) *ServicesAggregator {
	return &ServicesAggregator{
		EventEmitter: ee,
		current:      map[string]swarm.Service{},
	}
}

func (sa *ServicesAggregator) OnServices(services []swarm.Service) {
	sa.Lock()
	defer sa.Unlock()

	newServiceList := ServiceList{}

	newServiceMap := map[string]swarm.Service{}

	for _, s := range services {
		newServiceList = append(newServiceList, ServiceStatus{Name: s.Spec.Name, ID: s.ID})
		newServiceMap[s.ID] = s

		if _, found := sa.current[s.ID]; found {
			continue
		}

		sa.current[s.ID] = s

		sa.Emit(fmt.Sprintf("update/%s", s.ID), s)

	}

	for k := range sa.current {
		if _, found := newServiceMap[k]; !found {
			sa.Emit(fmt.Sprintf("delete/%s", k))
		}
	}

	sa.current = newServiceMap

	sort.Sort(newServiceList)

	if !reflect.DeepEqual(sa.serviceList, newServiceList) {
		sa.serviceList = newServiceList
		sa.Emit("list", newServiceList)
	}

}
