package logger

import "fmt"

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

type Logger struct {
	Level string
}

func New(level string) *Logger {
	return &Logger{Level: level}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Warn(msg string) {
	if l.Level != LevelError {
		fmt.Printf("WARN: %s\n", msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.Level == LevelDebug {
		fmt.Printf("DEBUG: %s\n", msg)
	}
}

func (l Logger) Error(msg string) {
	if l.Level != LevelInfo && l.Level != LevelWarn {
		fmt.Printf("ERROR: %s\n", msg)
	}
}
