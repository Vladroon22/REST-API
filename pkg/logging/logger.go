package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var ent *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLog() Logger {
	return Logger{ent}
}

type writeHook struct {
	Writer []io.Writer
	Level  []logrus.Level
}

func (hook *writeHook) Fire(ent *logrus.Entry) error {
	line, err := ent.String()
	if err != nil {
		return err
	}

	for _, wr := range hook.Writer {
		wr.Write([]byte(line))
	}
	return err
}

func (hook *writeHook) Levels() []logrus.Level {
	return hook.Level
}

func (l *Logger) GetLoggerWithField(k string, inter interface{}) Logger {
	return Logger{l.WithField(k, inter)}
}

func start_logging() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(fr *runtime.Frame) (function string, file string) {
			fileN := path.Base(fr.File)
			return fmt.Sprintf(fr.Function), fmt.Sprintf("%s:%d", fileN, fr.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0644)
	if err != nil {
		log.Fatal(err.Error())
	}

	file, err := os.OpenFile("logs/all_logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		log.Fatal(err.Error())
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writeHook{
		Writer: []io.Writer{file, os.Stdout},
		Level:  logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	ent = logrus.NewEntry(l)
}
