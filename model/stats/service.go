package stats

import "time"

var Service = NewContainerStats(time.Second * 120)
