package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogFileName = "./log.log"
	defaultLogFileSize = 100
	defaultBackups     = 5
	defaultMaxAge      = 7
)

var (
	std = newDefault()
)

// start default logger
func newDefault() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = customTimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(getLogFileWriter(defaultLogFileName, defaultLogFileSize, defaultBackups, defaultMaxAge)),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return true }),
	)
	return zap.New(core, zap.AddCaller())
}

// alias for zap logger
var (
	Debug = std.Debug
	Info  = std.Info
	Warn  = std.Warn
	Error = std.Error
	Fatal = std.Fatal
	Panic = std.Panic

	Print  = std.Debug
	Printf = Debugf
)

func Debugf(s string, args ...interface{}) {
	std.Debug(fmt.Sprintf(s, args...))
}

func Infof(s string, args ...interface{}) {
	std.Info(fmt.Sprintf(s, args...))
}

func Warnf(s string, args ...interface{}) {
	std.Warn(fmt.Sprintf(s, args...))
}

func Errorf(s string, args ...interface{}) {
	std.Error(fmt.Sprintf(s, args...))
}

func Fatalf(s string, args ...interface{}) {
	std.Error(fmt.Sprintf(s, args...))
}

func Panicf(s string, args ...interface{}) {
	std.Panic(fmt.Sprintf(s, args...))
}
