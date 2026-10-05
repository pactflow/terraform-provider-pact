package main

import (
	"sync"
	"testing"
)

func TestValidateEventsConcurrent(t *testing.T) {
	events := append([]string(nil), allowedEvents...) // snapshot before goroutines start
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		for _, e := range events {
			wg.Add(1)
			go func(e string) {
				defer wg.Done()
				if _, errs := validateEvents(e, "events.0"); len(errs) > 0 {
					t.Errorf("valid event %q rejected: %v", e, errs)
				}
			}(e)
		}
	}
	wg.Wait()
}
