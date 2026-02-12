package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type AppLog struct {
	*logrus.Logger
	logFile *os.File
}

func NewLog(namespace string) *AppLog {
	a := new(AppLog)

	locationFolder := "helpers/log/logfiles/"
	if err := os.MkdirAll(locationFolder, 0755); err != nil {
		a.Warnln("Error creating log folder:", err)
	}

	filename := fmt.Sprintf("%s_%s.log", namespace, time.Now().Format("20060102"))
	a.Logger = logrus.New()

	var err error
	if a.logFile, err = os.OpenFile(locationFolder+filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err == nil {
		o := io.MultiWriter(os.Stdout, a.logFile)
		a.SetOutput(o)
	} else {
		a.Warnln("saving log to file failed")
	}

	return a
}

func (a *AppLog) Close() {
	if a.logFile != nil {
		_ = a.logFile.Close()
	}
}
