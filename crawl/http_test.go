package crawl

import (
	"testing"
)

func TestNorm(t *testing.T) {
	value := norm(20, 10)
	t.Logf("value is %d", value)
}
