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

var (
	taskLogger         *Logger
	agentLogger        *Logger
	userHandlerLogger  *Logger
	adminHandlerLogger *Logger
)

type Logger struct {
	fp string
	*logrus.Logger
}

func New(fp string) *Logger {
	return &Logger{
		fp: filepath.Join("log", fp),
	}
}

func GetTaskLogger() *Logger {
	if taskLogger != nil {
		return taskLogger
	}
	return New(config.C.Log.Task)
}

func GetAgentLogger() *Logger {
	if agentLogger != nil {
		return agentLogger
	}
	return New(config.C.Log.Agent)
}

func GetUserHandlerLogger() *Logger {
	if userHandlerLogger != nil {
		return userHandlerLogger
	}
	return New(config.C.Log.UserHandler)
}

func GetAdminHandlerLogger() *Logger {
	if adminHandlerLogger != nil {
		return adminHandlerLogger
	}
	return New(config.C.Log.AdminHandler)
}

func MustInit() {
	taskLogger = GetTaskLogger().MustInit()
	agentLogger = GetAgentLogger().MustInit()
	userHandlerLogger = GetUserHandlerLogger().MustInit()
	adminHandlerLogger = GetAdminHandlerLogger().MustInit()
}

func (l *Logger) MustInit() *Logger {
	file, err := os.OpenFile(l.fp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("cannot open file: %s, error: %v", l.fp, err)
	}

	l.Logger = logrus.New()
	l.Out = io.MultiWriter(os.Stdout, file)

	lv, _ := logrus.ParseLevel(config.C.Log.Level)
	l.SetLevel(lv)

	return l
}
