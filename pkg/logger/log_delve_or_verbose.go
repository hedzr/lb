//go:build delve || verbose || hzstudio
// +build delve verbose hzstudio

package logger

import "log"

const LogValid bool = true

func Log(format string, v ...interface{}) { //nolint:goprintffuncname //no
	if realLogger != nil {
		realLogger.Printf(format, v...)
		return
	}
	// log.Skip(0).Infof(format, args...) // is there a `log` bug? so Skip(0) is a must-have rather than Skip(1), because stdLogger will detect how many frames should be skipped
	log.Printf(format, v...)
}

//

var realLogger Logger

func SetLogger(l Logger) (old Logger) {
	old = realLogger
	realLogger = l
	return
}

func SetLevel(lvl Level) {
	if realLogger != nil {
		realLogger.SetLevel(lvl)
	}
}

func GetLevel() (lvl Level) {
	if realLogger != nil {
		lvl = realLogger.GetLevel()
	}
	return
}

func Infof(format string, v ...interface{}) {
	if realLogger != nil {
		realLogger.Infof(format, v...)
		return
	}
	Log(format, v...)
}

func Warnf(format string, v ...interface{}) {
	if realLogger != nil {
		realLogger.Warnf(format, v...)
		return
	}
	Log(format, v...)
}

func Debugf(format string, v ...interface{}) {
	if realLogger != nil {
		realLogger.Debugf(format, v...)
		return
	}
	Log(format, v...)
}

func Tracef(format string, v ...interface{}) {
	if realLogger != nil {
		realLogger.Tracef(format, v...)
		return
	}
	Log(format, v...)
}

func Errorf(format string, v ...interface{}) {
	if realLogger != nil {
		realLogger.Errorf(format, v...)
		return
	}
	Log(format, v...)
}

func Printf(format string, v ...interface{}) {
	if realLogger != nil {
		realLogger.Printf(format, v...)
		return
	}
	Log(format, v...)
}
