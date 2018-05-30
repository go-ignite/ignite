package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-ignite/ignite/config"

	"github.com/sirupsen/logrus"
)

func init() {
	os.Mkdir("log", os.ModePerm)
}

func New(fp string) *logrus.Logger {
	// new logrus logger
	fp = filepath.Join("log", fp)
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("cannot open file: %s, error: %v", fp, err)
	}

	l := logrus.New()
	l.Out = io.MultiWriter(os.Stdout, file)

	// set log level
	lv, _ := logrus.ParseLevel(config.C.App.LogLevel)
	l.SetLevel(lv)

	return l
}
