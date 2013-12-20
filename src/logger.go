package collector

import (
/*. "./intf"
"fmt"*/
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

func GetLogger() *Logger {
	return loggerInstance
}

func NewLogger(config LoggerConfig) {
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

func (self *Logger) Info(log interface{}, args ...interface{}) {
	self.log4go.Info(log, args...)
}

func (self *Logger) Debug(log interface{}, args ...interface{}) {
	self.log4go.Debug(log, args...)
}

func (self *Logger) Warning(log interface{}, args ...interface{}) {
	self.log4go.Warn(log, args...)
}

func (self *Logger) Error(log interface{}, args ...interface{}) {
	self.log4go.Error(log, args...)
}

func (self *Logger) Critical(log interface{}, args ...interface{}) {
	self.log4go.Critical(log, args...)
}
