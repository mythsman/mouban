package crawl

import (
	"fmt"
	"mouban/internal/common"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	common.InitConfigBase()

	if err := ensureClientsInitialized(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "crawl test bootstrap failed: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
