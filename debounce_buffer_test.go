package debounce

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPerDeviceBufferingAndSmoothing(t *testing.T) {
	const numDevices = 4
	const flushAfter = 2 * time.Second
	const sendDelay = 500 * time.Millisecond

	var mu sync.Mutex
	sentTimestamps := make(map[string][]time.Time)

	mgr := NewManager(Config{
		FlushAfter:                2 * time.Second,
		MinIntervalBetweenFlushes: 200 * time.Millisecond,
		SendFunc: func(key string, batch []interface{}) {
			now := time.Now()
			mu.Lock()
			sentTimestamps[key] = append(sentTimestamps[key], now)
			mu.Unlock()
			fmt.Printf("[Send] %s: %d items\n", key, len(batch))
		},
	})

	start := time.Now()
	// Simulate random arrivals from 10 devices
	var wg sync.WaitGroup
	for i := 0; i < numDevices; i++ {
		deviceID := fmt.Sprintf("dev-%02d", i)
		wg.Add(1)

		go func(deviceID string) {
			defer wg.Done()
			for time.Since(start) < 5*time.Second {
				n := rand.Intn(10) + 1 // 1–10
				for j := 0; j < n; j++ {
					mgr.Add(deviceID, fmt.Sprintf("%s-fp-%d", deviceID, j))
				}
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			}
		}(deviceID)
	}

	wg.Wait()
	time.Sleep(3 * time.Second) // allow final flushes
}
