package common

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLog(t *testing.T) {
	logrus.Infolnln("test", "test", "test")
}
