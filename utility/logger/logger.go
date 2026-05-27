package logger

import "github.com/LandcLi/landc-go/log/facade"

var globalLogger facade.Logger

func InitLogger(config *facade.LogConfig) {
	if globalLogger == nil {
		globalLogger = facade.GetLoggerWithLogConfig(config)
	}
}

func GetLogger() facade.Logger { return globalLogger }
