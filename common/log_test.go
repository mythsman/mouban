package common

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLog(t *testing.T) {
	logrus.Infoln("test", "test", "test")
}
