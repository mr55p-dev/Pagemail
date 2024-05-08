package timer

import "time"

type cronTimer struct {
	now  time.Time
	T    <-chan time.Time
	done chan struct{}
}

func (t *cronTimer) Stop() {
	close(t.done)
}

func NewCronTimer(interval time.Duration, start time.Time) *cronTimer {
	timer := new(cronTimer)
	ch := make(chan time.Time)
	done := make(chan struct{})
	timer.T = ch
	timer.done = done
	go func() {
		// get the first tick
		now := time.Now()
		nextTick := start
		for nextTick.Before(now) {
			nextTick = nextTick.Add(interval)
		}

		for {
			select {
			case <-done:
				close(ch)
				return
			case t := <-time.After(nextTick.Sub(time.Now())):
				ch <- t
				nextTick = nextTick.Add(interval)
			}
		}
	}()

	return timer
}
