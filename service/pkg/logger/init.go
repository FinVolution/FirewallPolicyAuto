package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logFileName, logLevel string, fileSize, backups, age int) {
	// 配置日志文件的最大大小（以 MB 为单位）
	if fileSize <= 0 {
		fileSize = defaultLogFileSize
	}

	// 配置保留的日志文件数量
	if backups <= 0 {
		backups = defaultBackups
	}

	// 配置保留的日志天数
	if age <= 0 {
		age = defaultMaxAge
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = customTimeEncoder

	// 创建日志核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(getLogFileWriter(logFileName, fileSize, backups, age)),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= transLogLevel(logLevel) }),
	)

	// 添加文件名和行号到 Logger 中
	std = std.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core { return core }), zap.AddCaller(), zap.AddCallerSkip(1))
	defer std.Sync()
}

// 获取日志文件写入器
func getLogFileWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,    // 单个日志文件的最大大小（以 MB 为单位）
		MaxBackups: maxBackups, // 保留的日志文件数量
		MaxAge:     maxAge,
		LocalTime:  true,
	}

	return zapcore.AddSync(lumberjackLogger)
}

// 自定义时间格式化
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// 日志等级转换
func transLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "FATAL":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}
