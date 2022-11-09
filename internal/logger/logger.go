package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Info - пишет в лог успешно выполненные операции
func (l *Logger) Info(msg string) {
	l.Println("INFO: ", msg)
}

// Error - пишет в лог неудачно выполненные операции
func (l *Logger) Error(msg string) {
	l.Println("ERROR: ", msg)
}
