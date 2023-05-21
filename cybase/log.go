package cybase

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func ecszapInfoAndErrorLogger(appName string) (*zap.Logger, *zap.Logger) {
	var infoZapLogger *zap.Logger
	var errorZapLogger *zap.Logger
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	{
		infoCore := ecszap.NewCore(encoderConfig, os.Stdout, zap.InfoLevel)
		infoZapLogger = zap.New(infoCore, zap.AddCaller()).Named(appName)
	}
	{
		infoCore := ecszap.NewCore(encoderConfig, os.Stderr, zap.InfoLevel)
		errorZapLogger = zap.New(infoCore, zap.AddCaller()).Named(appName)
	}
	return infoZapLogger, errorZapLogger
}
