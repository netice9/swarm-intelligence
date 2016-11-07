package stats

import "time"

type Entry struct {
	Time   time.Time
	CPU    float64
	Memory int64
}

type Tracker struct {
	duration time.Duration
	entries  []Entry
}

func NewTracker(duration time.Duration) *Tracker {
	return &Tracker{duration: duration}
}

func (t *Tracker) Add(entry Entry) {

	defer func() {
		if len(t.entries) == 0 {
			return
		}
		lastTime := t.entries[len(t.entries)-1].Time
		minTime := lastTime.Add(-t.duration)
		for i, e := range t.entries {
			if e.Time == minTime || e.Time.After(minTime) {
				t.entries = t.entries[i:]
				return
			}
		}
		t.entries = []Entry{}
	}()

	for i, e := range t.entries {
		if e.Time.After(entry.Time) {
			t.entries = append(t.entries[:i], append([]Entry{entry}, t.entries[i:]...)...)
			return
		}
	}
	t.entries = append(t.entries, entry)
}

func (t *Tracker) Entries() []Entry {
	return t.entries
}
