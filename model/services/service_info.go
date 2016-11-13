package services

import (
	"reflect"
	"sort"

	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/model/stats"
)

type ServiceInfo struct {
	Service swarm.Service
	Tasks   map[string]TaskInfo
	Mem     uint64
	CPU     uint64
}

func NewServiceInfo(service swarm.Service) *ServiceInfo {
	return &ServiceInfo{
		Service: service,
		Tasks:   map[string]TaskInfo{},
	}
}

func (s *ServiceInfo) updateStats() {
	cumulative := stats.Entry{}
	for _, t := range s.Tasks {
		cumulative = cumulative.Add(stats.Service.LastStats(t.ContainerID))
	}
	s.CPU = cumulative.CPU
	s.Mem = cumulative.Memory
}

func (s *ServiceInfo) UpdateTasks(tasks []swarm.Task) bool {
	changed := false

	newTasks := map[string]TaskInfo{}

	for _, task := range tasks {
		current, found := s.Tasks[task.ID]
		currentTask := NewTaskInfo(task)
		newTasks[task.ID] = currentTask
		if !found {
			changed = true
			continue
		}
		if !reflect.DeepEqual(currentTask, current) {

			changed = true
		}
	}

	for id := range s.Tasks {
		if _, found := newTasks[id]; !found {
			changed = true
		}
	}
	s.Tasks = newTasks
	return changed
}

func (s *ServiceInfo) GetTasks() TaskInfoList {
	tasks := TaskInfoList{}
	for _, t := range s.Tasks {
		tasks = append(tasks, t)
	}
	sort.Sort(tasks)
	return tasks
}
