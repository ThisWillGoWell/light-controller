package log

import "go.uber.org/zap"

type Logger struct {
	zapLogger *zap.SugaredLogger
}

func (l *Logger) Warn(msg string, with ...interface{}) {

}

func (l *Logger) Info(msg string, with ...interface{}) {

}

func (l *Logger) Error(msg string, err error, with ...interface{}) {
}
