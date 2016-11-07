package services

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/event"
)

type ServiceStatus struct {
	Name string
	ID   string
}

type ServicesAggregator struct {
	sync.Mutex
	current map[string]*ServiceInfo
	event.EventEmitter
	serviceList ServiceList
}

func NewServicesAggregator(ee event.EventEmitter) *ServicesAggregator {
	return &ServicesAggregator{
		EventEmitter: ee,
		current:      map[string]*ServiceInfo{},
	}
}

func (sa *ServicesAggregator) ServiceList() []ServiceStatus {
	sa.Lock()
	defer sa.Unlock()
	return sa.serviceList
}

func (sa *ServicesAggregator) GetServiceInfo(serviceID string) *ServiceInfo {
	sa.Lock()
	defer sa.Unlock()
	return sa.current[serviceID]
}

func (sa *ServicesAggregator) OnTasks(tasks []swarm.Task) {

	sa.Lock()
	defer sa.Unlock()

	tasksByService := map[string][]swarm.Task{}
	for id := range sa.current {
		tasksByService[id] = []swarm.Task{}
	}

	for _, t := range tasks {
		_, found := sa.current[t.ServiceID]
		if found {
			tasksByService[t.ServiceID] = append(tasksByService[t.ServiceID], t)
			continue
		}
		log.Printf("Got task for not existing service %#v\n", t)
	}

	for id, s := range sa.current {
		updated := s.UpdateTasks(tasksByService[id])
		if updated {
			sa.Emit(fmt.Sprintf("update/%s", id), s)
		}
	}

}

func (sa *ServicesAggregator) OnServices(services []swarm.Service) {
	sa.Lock()
	defer sa.Unlock()

	newServiceList := ServiceList{}

	newServiceMap := map[string]*ServiceInfo{}

	for _, s := range services {
		newServiceList = append(newServiceList, ServiceStatus{Name: s.Spec.Name, ID: s.ID})

		if current, found := sa.current[s.ID]; found {

			if !reflect.DeepEqual(current.Service, s) {
				// newServiceMap[s.ID] = s
				current.Service = s
				sa.Emit(fmt.Sprintf("update/%s", s.ID), s)
			}

			newServiceMap[s.ID] = current

			continue
		}

		serviceInfo := NewServiceInfo(s)

		newServiceMap[s.ID] = serviceInfo

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
