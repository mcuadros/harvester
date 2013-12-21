package logger

import (
	"fmt"
	"os"
)

import "github.com/jarod/log4go"

type LoggerConfig struct {
	Level  string
	Format string
	File   string
}

type Logger struct {
	log4go log4go.Logger
}

var loggerInstance *Logger = nil

func init() {
	config := LoggerConfig{Format: "stdout", Level: "debug"}
	NewLogger(&config)
}

func GetLogger() *Logger {
	return loggerInstance
}

func NewLogger(config *LoggerConfig) {
	logger := new(Logger)

	level := log4go.WARNING
	if config.Level == "debug" {
		level = log4go.DEBUG
	} else if config.Level == "info" {
		level = log4go.INFO
	}

	logger.log4go = log4go.NewDefaultLogger(level)

	if config.Format == "log" {
		logger.log4go.AddFilter("log", level, log4go.NewFileLogWriter(config.File, true))
	} else {
		logger.log4go.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	loggerInstance = logger
}

/*
func (self *Logger) PrintWriterStats(elapsed int, writer Writer) {
	created, failed, transferred := writer.GetCounters()
	writer.ResetCounters()

	logFormat := "Created %d document(s), Failed %d times(s), %g"

	rate := float64(transferred) / 1000 / float64(elapsed)
	self.log4go.Info(fmt.Sprintf(logFormat, created, failed, rate))
}*/

func Info(log interface{}, args ...interface{}) {
	fmt.Println(log)
	loggerInstance.log4go.Info(log, args...)
}

func Debug(log interface{}, args ...interface{}) {
	fmt.Println(log)
	loggerInstance.log4go.Debug(log, args...)
}

func Warning(log interface{}, args ...interface{}) {
	fmt.Println(log)
	loggerInstance.log4go.Warn(log, args...)
}

func Error(log interface{}, args ...interface{}) {
	fmt.Println(log)
	loggerInstance.log4go.Error(log, args...)
}

func Critical(log interface{}, args ...interface{}) {
	fmt.Println(log)
	loggerInstance.log4go.Critical(log, args...)
	os.Exit(1)
}
