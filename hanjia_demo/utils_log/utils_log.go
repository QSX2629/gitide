package utils_log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"hanjia_demo/utils_viper"
)

// Logger 全局日志实例
var Logger *zap.Logger

// InitLogger 初始化日志（适配Viper配置）
func InitLogger() error {
	// 1. 获取最新配置（调用补全的GetConfig函数）
	cfg := utils_viper.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未初始化，请先调用utils_viper.InitConfig()")
	}
	logConfig := cfg.Log

	// 2. 创建日志目录（不存在则自动创建）
	logDir := filepath.Dir(logConfig.OutputPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 3. 构建日志核心
	core := buildLogCore(logConfig)

	// 4. 创建Logger实例
	options := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)}
	if utils_viper.IsDev() { // 调用utils_viper的IsDev函数
		options = append(options, zap.Development()) // 开发环境增强日志
	}
	Logger = zap.New(core, options...)
	zap.ReplaceGlobals(Logger) // 替换zap全局Logger

	Logger.Info("日志初始化成功", zap.String("log_level", logConfig.Level))
	return nil
}

// buildLogCore 构建日志核心
func buildLogCore(logConfig utils_viper.LogConfig) zapcore.Core {
	// 解析日志级别
	level := zapcore.InfoLevel
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	case "panic":
		level = zapcore.PanicLevel
	}

	// 编码器配置（统一格式）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 选择编码器（JSON/控制台）
	var encoder zapcore.Encoder
	if logConfig.Format == "json" || !utils_viper.IsDev() {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // 生产环境JSON格式
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 开发环境控制台格式
	}

	// 输出介质（文件+控制台）
	writers := []zapcore.WriteSyncer{}
	// 1. 文件输出（自动切割/压缩）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.OutputPath,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackup,
		MaxAge:     logConfig.MaxAge,
		Compress:   logConfig.Compress,
	})
	writers = append(writers, fileWriter)
	// 2. 开发环境控制台输出
	if utils_viper.IsDev() {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// 构建日志核心
	return zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), level)
}

// customTimeEncoder 自定义时间格式（带毫秒）
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}
