package gotest

import (
	"muslog"
	"os"
	"testing"
)

func Test_logInit(t *testing.T) {
	muslog.InitLog(muslog.DEBUG, "", muslog.LevelTrace)
	muslog.Trace("log should be print at stdout")
}

func Test_logFile(t *testing.T) {
	logFile, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	muslog.SetOutput(logFile)
	muslog.Trace("log should be at test.log")
}

func Test_trace(t *testing.T) {
	muslog.SetOutput(os.Stdout)
	muslog.Trace("log should be print at stdout")
}

func Test_Info(t *testing.T) {
	muslog.SetLevel(muslog.LevelError)
	muslog.Info("log should not be print at stdout")
}
