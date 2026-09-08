// Copyright (c) Alibaba, Inc. and its affiliates.

package common

import (
	"log"
	"os"
	"strings"
	"sync"
)

// Logger dashscope logger.
type Logger struct {
	info  *log.Logger
	debug *log.Logger
	warn  *log.Logger
	err   *log.Logger
	level int // 0 silent, 1 info, 2 debug
	mu    sync.Mutex
}

const (
	levelSilent = 0
	levelInfo   = 1
	levelDebug  = 2
)

// Log default logger.
var Log = newLogger()

func newLogger() *Logger {
	l := &Logger{
		info:  log.New(os.Stderr, "", log.LstdFlags),
		debug: log.New(os.Stderr, "", log.LstdFlags),
		warn:  log.New(os.Stderr, "", log.LstdFlags),
		err:   log.New(os.Stderr, "", log.LstdFlags),
		level: levelSilent,
	}
	level := os.Getenv(DASHSCOPE_LOGGING_LEVEL_ENV)
	if level != "" {
		if level != "info" && level != "debug" {
			level = "info"
		}
		if level == "debug" {
			l.level = levelDebug
		} else {
			l.level = levelInfo
		}
	}
	return l
}

// EnableLogging enable dashscope log for debugger.
func EnableLogging() {
	Log.mu.Lock()
	defer Log.mu.Unlock()
	level := os.Getenv(DASHSCOPE_LOGGING_LEVEL_ENV)
	if level == "" {
		return
	}
	if level != "info" && level != "debug" {
		level = "info"
	}
	if strings.EqualFold(level, "debug") {
		Log.level = levelDebug
	} else {
		Log.level = levelInfo
	}
}

func (l *Logger) Info(format string, args ...any) {
	if l.level >= levelInfo {
		l.info.Printf("[INFO] "+format, args...)
	}
}

func (l *Logger) Debug(format string, args ...any) {
	if l.level >= levelDebug {
		l.debug.Printf("[DEBUG] "+format, args...)
	}
}

func (l *Logger) Warn(format string, args ...any) {
	if l.level >= levelInfo {
		l.warn.Printf("[WARN] "+format, args...)
	}
}

func (l *Logger) Error(format string, args ...any) {
	l.err.Printf("[ERROR] "+format, args...)
}

func init() {
	EnableLogging()
}
