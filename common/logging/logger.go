package logging

import (
	"fmt"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	logger = log.StandardLogger()
)

func init() {
	logger.SetLevel(log.Level(NewLevelFromEnv()))
}

func generateFields() log.Fields {
	pc, _, line, _ := runtime.Caller(2)
	class := runtime.FuncForPC(pc).Name()
	where := fmt.Sprintf("%v:%v", class, line)
	return log.Fields{
		"where": where,
		"when":  time.Now().Format("01-02:15-04-05"),
	}
}

func Panic(format string, args ...interface{}) {
	logger.WithFields(generateFields()).Panicf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.WithFields(generateFields()).Fatalf(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.WithFields(generateFields()).Errorf(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.WithFields(generateFields()).Warnf(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.WithFields(generateFields()).Infof(format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.WithFields(generateFields()).Debugf(format, args...)
}

func Trace(format string, args ...interface{}) {
	logger.WithFields(generateFields()).Tracef(format, args...)
}
