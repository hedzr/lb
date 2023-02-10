//go:build !delve && !verbose && !hzstudio
// +build !delve,!verbose,!hzstudio

package logger

const LogValid bool = false //nolint:gochecknoglobals //i know that

func Log(format string, v ...interface{}) { //nolint:goprintffuncname //no

}

//

func SetLogger(l Logger) (old Logger) {
	return
}

func SetLevel(lvl Level) {
}

func GetLevel() (lvl Level) {
	return WarnLevel
}

func Infof(format string, v ...interface{}) {

}

func Warnf(format string, v ...interface{}) {

}

func Debugf(format string, v ...interface{}) {

}

func Tracef(format string, v ...interface{}) {

}

func Errorf(format string, v ...interface{}) {

}

func Fatalf(format string, v ...interface{}) {

}

func Panicf(format string, v ...interface{}) {

}

func Printf(format string, v ...interface{}) {

}
