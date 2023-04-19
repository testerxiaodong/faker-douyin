package log

import (
	"faker-douyin/internal/app/config"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var AppLogger *zap.Logger

func NewLogger(conf *config.Config) *zap.Logger {
	_, err := os.Stat(conf.Server.LogDir)
	if err != nil {
		if os.IsNotExist(err) && !config.IsDev() {
			err := os.MkdirAll(conf.Server.LogDir, os.ModePerm)
			if err != nil {
				panic("mkdir failed![%v]")
			}
		}
	}

	var core zapcore.Core

	if config.IsDev() {
		// 开发环境，使用开发环境的编码器
		core = zapcore.NewCore(getDevEncoder(), os.Stdout, getLogLevel(conf.Log.Levels.App))
	} else {
		// 生产环境，使用生产环境的编码器
		core = zapcore.NewCore(getProdEncoder(), getWriter(conf.Server.LogDir, conf.Log.FileName, conf.Log.MaxSize, conf.Log.MaxAge, conf.Log.Compress), zap.DebugLevel)
	}

	// 传入 log.AddCaller() 显示打日志点的文件名和行数
	AppLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))

	return AppLogger
}

// getWriter 自定义Writer,分割日志
func getWriter(logDir string, fileName string, maxSize int, maxAge int, compress bool) zapcore.WriteSyncer {
	rotatingLogger := &lumberjack.Logger{
		Filename: filepath.Join(logDir, fileName),
		MaxSize:  maxSize,
		MaxAge:   maxAge,
		Compress: compress,
	}
	return zapcore.AddSync(rotatingLogger)
}

// getProdEncoder 自定义日志编码器
func getProdEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 生产环境，使用Json格式日志
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getDevEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 开发环境，使用控制台行输出编码器
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		panic("log level error")
	}
}
