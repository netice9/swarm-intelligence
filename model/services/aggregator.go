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
	Name   string
	ID     string
	CPU    uint64
	Memory uint64
}

type ServicesAggregator struct {
	sync.Mutex
	current map[string]*ServiceInfo
	event.EventEmitter
	serviceList ServiceList
}

var Aggregator = NewServicesAggregator()

func NewServicesAggregator() *ServicesAggregator {

	services := &ServicesAggregator{
		EventEmitter: event.NewEmitterAdapter(),
		current:      map[string]*ServiceInfo{},
	}
	event.Services.AddListener("update", services.OnServices)
	event.Tasks.AddListener("update", services.OnTasks)
	event.Time.AddListener("1sec", services.OnTimer)

	return services

}

func (sa *ServicesAggregator) OnServiceList(fn func(ServiceList)) {
	sa.Lock()
	defer sa.Unlock()
	list := sa.serviceList
	go fn(list)
	sa.AddListener("list", fn)
}

func (sa *ServicesAggregator) RemoveServiceListListener(fn func(ServiceList)) {
	sa.RemoveListener("list", fn)
}

func (sa *ServicesAggregator) OnServiceInfo(serviceID string, fn func(*ServiceInfo)) {
	sa.Lock()
	defer sa.Unlock()
	info := sa.current[serviceID]
	go fn(info)
	sa.AddListener(fmt.Sprintf("update/%s", serviceID), fn)
}

func (sa *ServicesAggregator) RemoveServiceInfoListener(serviceID string, fn func(*ServiceInfo)) {
	sa.RemoveListener(fmt.Sprintf("update/%s", serviceID), fn)
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
			s.updateStats()
			sa.Emit(fmt.Sprintf("update/%s", id), s)
		}

	}

}

func (sa *ServicesAggregator) OnTimer() {
	sa.Lock()
	defer sa.Unlock()
	defer sa.updateServiceList()
	for id, s := range sa.current {
		s.updateStats()
		sa.Emit(fmt.Sprintf("update/%s", id), s)
	}

}

func (sa *ServicesAggregator) OnServices(services []swarm.Service) {
	sa.Lock()
	defer sa.Unlock()
	defer sa.updateServiceList()

	newServiceMap := map[string]*ServiceInfo{}

	for _, s := range services {

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

}

func (sa *ServicesAggregator) updateServiceList() {
	newServiceList := ServiceList{}

	for _, si := range sa.current {
		newServiceList = append(newServiceList, si.Status())
	}

	sort.Sort(newServiceList)

	if !reflect.DeepEqual(sa.serviceList, newServiceList) {
		sa.serviceList = newServiceList
		sa.Emit("list", newServiceList)
	}

	sa.serviceList = newServiceList

}
