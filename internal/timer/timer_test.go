package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	now := time.Now()
	start := now.Add(time.Second * 2)
	timer := NewCronTimer(time.Second*1, start)
	<-time.After(time.Second)
	assert.Empty(t, timer.T)
	for i := 0; i < 3; i++ {
		select {
		case recv := <-timer.T:
			recv = recv.Round(time.Second)
			now := time.Now().Round(time.Second)
			assert.Equal(t, recv, now)
		case <-time.After(time.Second * 2):
			assert.Fail(t, "timeout")
		}
	}
	timer.Stop()
	assert.Empty(t, timer.T)
}
