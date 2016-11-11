package stats

import "sync"

type ContainersTracker struct {
	sync.Mutex
	trackers map[string]*Tracker
}
