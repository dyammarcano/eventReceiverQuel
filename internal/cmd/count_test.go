package cmd

import (
	"testing"
	"time"
)

func TestNewCounterSvc(t *testing.T) {
	c := NewCounterSvc()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {

		select {
		case <-ticker.C:
			c.CountEvent()
		}
	}()
	c.Start()
}
