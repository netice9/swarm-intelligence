package services

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

type TaskInfo struct {
	ID           string
	State        string
	DesiredState string
	Slot         int
	NodeID       string
	ContainerID  string
	CreatedAt    time.Time
}

func NewTaskInfo(t swarm.Task) TaskInfo {
	return TaskInfo{
		ID:           t.ID,
		State:        string(t.Status.State),
		DesiredState: string(t.DesiredState),
		Slot:         t.Slot,
		NodeID:       t.NodeID,
		ContainerID:  t.Status.ContainerStatus.ContainerID,
		CreatedAt:    t.CreatedAt,
	}
}

type TaskInfoList []TaskInfo

func (t TaskInfoList) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t TaskInfoList) Less(i, j int) bool {
	if t[i].CreatedAt == t[j].CreatedAt {
		return t[i].ID < t[j].ID
	}
	return !t[i].CreatedAt.Before(t[j].CreatedAt)
}

// Swap swaps the elements with indexes i and j.
func (t TaskInfoList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
