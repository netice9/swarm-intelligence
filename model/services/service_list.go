package services

type ServiceList []ServiceStatus

func (sl ServiceList) Len() int {
	return len(sl)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (sl ServiceList) Less(i, j int) bool {
	if sl[i].CreatedAt == sl[j].CreatedAt {
		if sl[i].Name == sl[j].Name {
			return sl[i].ID < sl[j].ID
		}
		return sl[i].Name < sl[j].Name
	}
	return sl[i].CreatedAt.Before(sl[j].CreatedAt)
}

// Swap swaps the elements with indexes i and j.
func (sl ServiceList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}
