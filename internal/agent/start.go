package agent

import "sync"

var startOnce sync.Once

// Start initializes all background agents exactly once.
func Start() {
	startOnce.Do(func() {
		startCounterAgent()
		startFallbackAgent()
		startLatestAgent()
		startUserAgent()
		startItemAgent()
	})
}
