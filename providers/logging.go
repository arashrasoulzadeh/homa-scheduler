package providers

import "go.uber.org/zap"

func LoggingSugar() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}
