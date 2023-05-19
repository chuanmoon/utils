package mlog

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func EcszapLogger(appName string) *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	infoCore := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	infoZapLogger := zap.New(infoCore, zap.AddCaller())
	infoZapLogger = infoZapLogger.Named(appName)
	return infoZapLogger
}

func EcszapInfoAndErrorLogger(appName string) (*zap.Logger, *zap.Logger) {
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
