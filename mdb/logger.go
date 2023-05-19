package mdb

import "go.uber.org/zap"

// InfoLogger 适配zap logger
type InfoLogger struct {
	zapLogger *zap.Logger
}

// WarnLoggor 适配zap logger
type WarnLoggor struct {
	zapLogger *zap.Logger
}

// ErrorLogger 适配zap logger
type ErrorLogger struct {
	zapLogger *zap.Logger
}

func loggerFields(data map[string]string) []zap.Field {
	var fields []zap.Field
	for key, val := range data {
		fields = append(fields, zap.String(key, val))
	}
	return fields
}

func (l *ErrorLogger) Log(msg string, data map[string]string) {
	if l == nil || l.zapLogger == nil {
		return
	}
	l.zapLogger.Error(msg, loggerFields(data)...)
}

func (l *InfoLogger) Log(msg string, data map[string]string) {
	if l == nil || l.zapLogger == nil {
		return
	}
	l.zapLogger.Info(msg, loggerFields(data)...)
}

func (l *WarnLoggor) Log(msg string, data map[string]string) {
	if l == nil || l.zapLogger == nil {
		return
	}
	l.zapLogger.Warn(msg, loggerFields(data)...)
}
