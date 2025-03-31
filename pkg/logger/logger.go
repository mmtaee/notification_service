package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type CustomLogger struct {
	ctx     context.Context
	cancel  context.CancelFunc
	logChan chan string
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Green   = "\033[32m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
)

var Logger *CustomLogger

func Init(bufferSize int) {
	log.SetFlags(0)
	ctx, cancel := context.WithCancel(context.Background())
	Logger = &CustomLogger{
		ctx:     ctx,
		cancel:  cancel,
		logChan: make(chan string, bufferSize),
	}
	go Logger.processLogs()
}

func (c *CustomLogger) processLogs() {
	for {
		select {
		case msg, ok := <-c.logChan:
			if !ok {
				return
			}
			log.Println(msg)
		case <-c.ctx.Done():
			close(c.logChan)
			return
		}
	}
}

func Close() {
	Logger.cancel()
}

func (c *CustomLogger) logMessage(level, color string, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	msg = fmt.Sprintf(
		"%s[%s] [%s] %s %s",
		color,
		level,
		time.Now().Format("2006-01-02 15:04:05"),
		msg,
		Reset,
	)
	c.logChan <- msg
}

func Info(msg string, args ...interface{}) {
	Logger.logMessage("INFO", Green, msg, args...)
}

func Success(msg string, args ...interface{}) {
	Logger.logMessage("INFO", Blue, msg, args...)
}

func Warning(msg string, args ...interface{}) {
	Logger.logMessage("WARNING", Yellow, msg, args...)
}

func Error(msg string, args ...interface{}) {
	Logger.logMessage("ERROR", Red, msg, args...)
}

func Critical(msg string, args ...interface{}) {
	msg = fmt.Sprintf(
		"%s[CRITICAL] [%s] %s %s",
		Magenta,
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(msg, args...),
		Reset,
	)
	log.Println(msg)
	os.Exit(1)
}
