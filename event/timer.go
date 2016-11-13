package event

import "time"

var Time = NewEmitterAdapter()

func StartTimerEvents() {

	count := 0
	go func() {
		for range time.Tick(time.Second) {
			Time.Emit("1sec")
			if count == 4 {
				Time.Emit("5sec")
				count = 0
				continue
			}
			count++
		}
	}()

}
