package logger

import (
	"go.uber.org/zap"
	"sync"
)

var logger *zap.Logger
var loggerOnce sync.Once

// 初始化 zap logger
func newLogger() (*zap.Logger, error) {
	// 可以选择生产环境配置（JSON 格式）或开发环境配置（带颜色控制台输出）
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func GetLogger() *zap.Logger {
	loggerOnce.Do(func() {
		logger, _ = newLogger()
	})
	return logger
}
