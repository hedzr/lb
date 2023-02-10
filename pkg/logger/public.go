// package logger is a minimal logger weapper
package logger

type Logger interface {
	Printf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Tracef(format string, v ...interface{})

	Errorf(format string, v ...interface{})

	SetLevel(lvl Level)
	GetLevel() Level
}
