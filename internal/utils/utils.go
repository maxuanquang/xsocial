package utils

import (
	"github.com/maxuanquang/social-network/configs"
	"go.uber.org/zap"
)

func NewLogger(cfg *configs.LoggerConfig) (*zap.Logger, error) {
	loggerCfg := zap.NewDevelopmentConfig()

	switch cfg.Level {
	case "debug":
		loggerCfg.Level.SetLevel(zap.DebugLevel)
	case "info":
		loggerCfg.Level.SetLevel(zap.InfoLevel)
	case "warn":
		loggerCfg.Level.SetLevel(zap.WarnLevel)
	case "error":
		loggerCfg.Level.SetLevel(zap.ErrorLevel)
	case "dpanic":
		loggerCfg.Level.SetLevel(zap.DPanicLevel)
	case "panic":
		loggerCfg.Level.SetLevel(zap.PanicLevel)
	case "fatal":
		loggerCfg.Level.SetLevel(zap.FatalLevel)
	}

	logger := zap.Must(loggerCfg.Build())
	return logger, nil
}
