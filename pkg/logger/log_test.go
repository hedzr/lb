//go:// build delve || verbose || hzstudio
// + // build delve verbose hzstudio

package logger_test

import (
	"testing"

	"github.com/hedzr/lb/pkg/logger"
)

func TestFLog(t *testing.T) {
	// config := log.NewLoggerConfigWith(true, "logrus", "trace")
	// logger := logrus.NewWithConfig(config)
	logger.Printf("hello")
	logger.Infof("hello info")
	logger.Warnf("hello warn")
	logger.Errorf("hello error")
	logger.Debugf("hello debug")
	logger.Tracef("hello trace")

	logger.Log("but again")
}
