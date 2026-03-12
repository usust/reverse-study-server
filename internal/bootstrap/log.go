package bootstrap

import (
	"strings"

	"github.com/usust/goUtils/logger"
)

func InitLogger() error {
	level := logger.LOG_LEVEL_DEBUG
	switch strings.ToLower(strings.TrimSpace(GlobalConfig.ZapLog.LogLevel)) {
	case "error":
		level = logger.LOG_LEVEL_ERROR
	case "warn", "warning":
		level = logger.LOG_LEVEL_WARN
	case "info":
		level = logger.LOG_LEVEL_INFO
	case "debug", "":
		level = logger.LOG_LEVEL_DEBUG
	default:
		level = logger.LOG_LEVEL_DEBUG
	}

	if err := logger.InitZapCore(nil,
		logger.ZapWithLogDir(GlobalConfig.ZapLog.LogDir),
		logger.ZapWithLevel(level),
	); err != nil {
		return err
	}

	// goUtils 初始化后从 zap 全局实例获取 logger。
	ZapLogger = logger.Logger
	SugaredLogger = logger.SugaredLogger
	SugaredLogger.Info("日志初始化成功")
	return nil
}
