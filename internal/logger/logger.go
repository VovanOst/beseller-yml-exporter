package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Level представляет уровень логирования
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger представляет простой логгер
type Logger struct {
	level  Level
	logger *log.Logger
}

// New создаёт новый логгер с указанным уровнем
func New(levelStr string) *Logger {
	level := parseLevel(levelStr)
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Debug логирует сообщение уровня DEBUG
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.log("DEBUG", msg, args...)
	}
}

// Info логирует сообщение уровня INFO
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log("INFO", msg, args...)
	}
}

// Warn логирует сообщение уровня WARN
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.log("WARN", msg, args...)
	}
}

// Error логирует сообщение уровня ERROR
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LevelError {
		l.log("ERROR", msg, args...)
	}
}

// Fatal логирует сообщение уровня ERROR и завершает программу
func (l *Logger) Fatal(args ...interface{}) {
	l.Error(fmt.Sprint(args...))
	os.Exit(1)
}

// log форматирует и выводит лог сообщение
func (l *Logger) log(level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Форматирование дополнительных аргументов
	var extra string
	if len(args) > 0 {
		extra = fmt.Sprintf(" %v", args)
	}

	l.logger.Printf("%s %s %s%s", timestamp, level, msg, extra)
}

// parseLevel парсит строку уровня логирования
func parseLevel(levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}
