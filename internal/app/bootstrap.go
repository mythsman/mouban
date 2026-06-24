package app

import (
	"fmt"
	"mouban/internal/agent"
	"mouban/internal/common"
	"mouban/internal/crawl"
)

// Bootstrap initializes runtime components in dependency order:
// common(config/logger/db) -> crawl(http clients) -> agent(background workers).
func Bootstrap() error {
	if err := bootstrapCommon(); err != nil {
		return fmt.Errorf("common bootstrap failed: %w", err)
	}

	if err := crawl.Bootstrap(); err != nil {
		return fmt.Errorf("crawl bootstrap failed: %w", err)
	}

	agent.Start()
	return nil
}

func bootstrapCommon() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	common.Bootstrap()
	return nil
}
