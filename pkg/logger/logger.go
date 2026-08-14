package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger // 日志记录器

func Init(level string) {
	var cfg zap.Config // 配置
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()                                 // 开发环境配置
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 日志级别使用颜色编码
	} else {
		cfg = zap.NewProductionConfig()                           // 生产环境配置
		cfg.EncoderConfig.TimeKey = "timestamp"                   // 时间键
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式化
	}

	// 可设置日志级别
	lv := zapcore.Level(0) // debug级别
	switch level {
	case "info":
		lv = zapcore.InfoLevel // 信息级别
	case "warn":
		lv = zapcore.WarnLevel // 警告级别
	case "error":
		lv = zapcore.ErrorLevel // 错误级别
	}
	cfg.Level.SetLevel(lv) // 设置日志级别

	var err error
	Log, err = cfg.Build() // 构建日志记录器
	if err != nil {
		panic(err) // 如果构建失败，则panic
	}
	zap.ReplaceGlobals(Log) // 替换 zap 全局 logger，使 zap.L() 可用
}

// 提供简写函数
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...) // 记录信息级别日志
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...) // 记录错误级别日志
}

func Sync() {
	_ = Log.Sync() // 同步日志记录器
}
