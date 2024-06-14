package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type StdLogger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
}

type stdLogger struct {
	scope string
}

func Std(scope string) StdLogger {
	return &stdLogger{scope}
}

func (l *stdLogger) log(s string) {
	New().
		WithOptions(
			zap.AddCallerSkip(2),
		).
		With(
			zap.String("source", "logger.StdLogger"),
			zap.String("scope", l.scope),
		).
		Info(s)
}

func (l *stdLogger) Print(v ...any) {
	l.log(fmt.Sprint(v...))
}

func (l *stdLogger) Printf(s string, v ...any) {
	l.log(fmt.Sprintf(s, v...))
}

func (l *stdLogger) Println(v ...any) {
	l.log(fmt.Sprint(v...))
}
