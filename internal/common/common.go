package common

import "sync"

var bootstrapOnce sync.Once

// Bootstrap initializes config, logger and database exactly once.
// Call this explicitly from main entrypoints or integration tests.
func Bootstrap() {
	bootstrapOnce.Do(func() {
		InitConfig()
		InitLogger()
		InitDatabase()
	})
}
