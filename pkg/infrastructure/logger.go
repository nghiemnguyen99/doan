package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	gorm "gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	return log
}

// Logger based on logrus, but compatible with gorm
type GormLogger struct {
	logger *logrus.Entry
}

func NewGormLogger(logger *logrus.Logger) GormLogger {
	return GormLogger{
		logger.WithField("service", "database"),
	}
}

// We ignore this setting, because the log level is already decided by logrus
func (logger GormLogger) LogMode(gormLog.LogLevel) gormLog.Interface {
	return logger
}

func (logger GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.logger.WithContext(ctx).Infof(msg, data)
}

func (logger GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.logger.WithContext(ctx).Warnf(msg, data)
}

func (logger GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.logger.WithContext(ctx).Errorf(msg, data)
}

func (logger GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	sql, rows := fc()
	duration := time.Since(begin)
	logEntry := logger.logger.
		WithContext(ctx).
		WithField("duration", duration.String()).
		WithField("sql", sql).
		WithField("rows", rows)

	if err == nil {
		logEntry.Info("Performed SQL Query")
	} else {
		logEntry = logEntry.WithField("error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logEntry.Info("Performed SQL Query")
		} else {
			logEntry.Error("SQL Query failed")
		}
	}
}
