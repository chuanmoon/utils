package mdb

import (
	"time"

	"github.com/yinshuwei/osm/v2"
	"go.uber.org/zap"
)

func Osm(driverName, dataSource string, infoLogger, errorLogger *zap.Logger) *osm.Osm {
	o, err := osm.New(driverName, dataSource, Options(infoLogger, errorLogger, 100*time.Millisecond))
	if err != nil {
		errorLogger.Error("osm new", zap.Error(err), zap.String("driver_name", driverName), zap.String("data_source", dataSource))
		panic(err)
	}
	return o
}

func Options(infoLogger, errorLogger *zap.Logger, slowLogDuration time.Duration) osm.Options {
	return osm.Options{
		MaxIdleConns:    50,
		MaxOpenConns:    100,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		WarnLogger:      &WarnLoggor{errorLogger},  // Logger
		ErrorLogger:     &ErrorLogger{errorLogger}, // Logger
		InfoLogger:      &InfoLogger{infoLogger},   // Logger
		ShowSQL:         true,                      // bool
		SlowLogDuration: slowLogDuration,           // time.Duration
	}
}
