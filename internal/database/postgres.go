package database

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CustomGormLogger bridges GORM logs to the application's logging system Zerolog
type CustomGormLogger struct{}

// LogMode sets the log level for the custom logger
func (l *CustomGormLogger) LogMode(level logger.LogLevel) logger.Interface { return l }
func (l *CustomGormLogger) Info(ctx context.Context, msg string, data ...any) {
	log.Ctx(ctx).Info().Msgf(msg, data...)
}
func (l *CustomGormLogger) Warn(ctx context.Context, msg string, data ...any) {
	log.Ctx(ctx).Warn().Msgf(msg, data...)
}
func (l *CustomGormLogger) Error(ctx context.Context, msg string, data ...any) {
	log.Ctx(ctx).Error().Msgf(msg, data...)
}
func (l *CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// pull logger from context which contains request_id (and request-scoped fields)
	lgr := log.Ctx(ctx)

	if err != nil && err != gorm.ErrRecordNotFound {
		lgr.Error().
			Err(err).
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Msg(sql)
	} else if elapsed > 200*time.Millisecond {
		lgr.Warn().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("perf", "SLOW_QUERY").
			Msgf("%s", sql)
	} else {
		// log all queries in development
		lgr.Debug().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Msg(sql)
	}
}

func Connect(dsn string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		AllowGlobalUpdate: false, // safety: prevent global updates/deletes without a WHERE clause
		PrepareStmt:       true,  // performance: cache prepared statements
	}

	gormConfig.Logger = &CustomGormLogger{}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)                  // Max active connections
	sqlDB.SetMaxIdleConns(10)                  // Max idle connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Reuse connections for an hour
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Close idle connections after 30 minutes

	return db, nil
}
