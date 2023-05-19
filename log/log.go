package log

import "go.uber.org/zap"

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	sugar = logger.Sugar()
}

func Sugar() *zap.SugaredLogger {
	return sugar
}

func Fatal(args ...any) {
	sugar.Fatal(args...)
}

func Fatalf(tmpl string, args ...any) {
	sugar.Fatalf(tmpl, args...)
}

func Info(args ...any) {
	sugar.Info(args...)
}

func Infof(tmpl string, args ...any) {
	sugar.Infof(tmpl, args...)
}

func Error(args ...any) {
	sugar.Error(args...)
}

func Errorf(tmpl string, args ...any) {
	sugar.Errorf(tmpl, args...)
}

func Clear() {
	logger.Sync()
}
