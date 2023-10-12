package common

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {
	logrus.Infoln("test", "test", "test")
}
