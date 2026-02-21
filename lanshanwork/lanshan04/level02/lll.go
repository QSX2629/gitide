package main

import (
	"fmt"
	"os"
	"time"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "Info"
	LogLevelWarn  LogLevel = "Warn"
	LogLevelError LogLevel = "Error"
	LogLevelFatal LogLevel = "Fatal"
)

type Logger struct {
	file     *os.File
	MinLevel LogLevel
}

func NewLogger(filename string, minLevel LogLevel) (*Logger, error) {
	file, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644)
	if err != nil {
		return nil, err
	}
	if minLevel == "" {
		minLevel = LogLevelInfo
	}
	return &Logger{
		file:     file,
		MinLevel: minLevel,
	}, nil
}
func (l *Logger) ShouldLog(level LogLevel) bool {
	LevelOrder := map[LogLevel]int{
		LogLevelInfo:  0,
		LogLevelWarn:  1,
		LogLevelError: 2,
		LogLevelFatal: 3,
	}
	return LevelOrder[level] >= LevelOrder[l.MinLevel]
}
func (l *Logger) LogWrite(level LogLevel, message string) error {
	if level < l.MinLevel {
		return nil
	}
	now := time.Now()
	timestamp := now.Format("2006-01-02 15:04:05")
	unixtime := now.Unix()
	nowtime := fmt.Sprintf("[%v] [%v] [%v] [%v]", timestamp, unixtime, level, message)
	_, err := l.file.WriteString(nowtime)
	return err
}
func (l *Logger) Info(message string) error {
	return l.LogWrite(LogLevelInfo, message)
}
func (l *Logger) Warn(message string) error {
	return l.LogWrite(LogLevelWarn, message)
}
func (l *Logger) Error(message string) error {
	return l.LogWrite(LogLevelError, message)
}

func (l *Logger) Fatal(message string) {
	_ = l.LogWrite(LogLevelFatal, message)
	l.Close()
	os.Exit(1)
}
func (l *Logger) Close() error {
	return l.file.Close()
}
func main() {
	logger, err := NewLogger("lll", LogLevelInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Close()
	logger.Info("程序启动成功（Info）")
	logger.Warn("Warn")
	logger.Error("Error")
	fmt.Println("OK")
}
