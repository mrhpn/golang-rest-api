package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CustomGormLogger bridges GORM logs to the application's logging system Zerolog
type CustomGormLogger struct{}

// LogMode sets the log level for the custom logger
// The log level is ignored because logging is delegated to Zerolog.
func (l *CustomGormLogger) LogMode(_ logger.LogLevel) logger.Interface { return l }

// Info logs informational messages from GORM.
func (l *CustomGormLogger) Info(ctx context.Context, msg string, data ...any) {
	log.Ctx(ctx).Info().Msgf(msg, data...)
}

// Warn logs warning messages from GORM.
func (l *CustomGormLogger) Warn(ctx context.Context, msg string, data ...any) {
	log.Ctx(ctx).Warn().Msgf(msg, data...)
}

// Error logs error messages from GORM.
func (l *CustomGormLogger) Error(ctx context.Context, msg string, data ...any) {
	log.Ctx(ctx).Error().Msgf(msg, data...)
}

// Trace logs SQL execution details including elapsed time and errors.
func (l *CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// pull logger from context which contains request_id (and request-scoped fields)
	lgr := log.Ctx(ctx)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

// Connect establishes a database connection with retry logic and configurable pool settings
func Connect(dsn string, dbCfg *config.DBConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		AllowGlobalUpdate: false, // safety: prevent global updates/deletes without a WHERE clause
		PrepareStmt:       true,  // performance: cache prepared statements
	}

	gormConfig.Logger = &CustomGormLogger{}

	var db *gorm.DB
	var err error

	// Retry connection with exponential backoff
	maxAttempts := dbCfg.RetryAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3
	}

	retryDelay := time.Duration(dbCfg.RetryDelaySecond) * time.Second
	if retryDelay <= 0 {
		retryDelay = 2 * time.Second
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}

		if attempt < maxAttempts {
			backoff := retryDelay * time.Duration(attempt)
			log.Warn().
				Err(err).
				Int("attempt", attempt).
				Int("max_attempts", maxAttempts).
				Dur("backoff", backoff).
				Msg("database connection failed, retrying...")
			time.Sleep(backoff)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	// Set max open connections
	maxOpenConns := dbCfg.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = 25
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)

	// Set max idle connections
	maxIdleConns := dbCfg.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// Set connection max life time
	connMaxLifetime := time.Duration(dbCfg.ConnMaxLifetimeMinute) * time.Minute
	if connMaxLifetime <= 0 {
		connMaxLifetime = time.Hour
	}
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// Set connection max idle time
	connMaxIdleTime := time.Duration(dbCfg.ConnMaxIdleTimeMinute) * time.Minute
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = 30 * time.Minute
	}
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Int("max_open_conns", maxOpenConns).
		Int("max_idle_conns", maxIdleConns).
		Dur("conn_max_lifetime", connMaxLifetime).
		Dur("conn_max_idle_time", connMaxIdleTime).
		Msg("database connection established successfully")

	// Register callback to automatically apply query timeout to all queries
	registerQueryTimeoutCallback(db, dbCfg)

	return db, nil
}

// registerQueryTimeoutCallback registers a GORM callback that automatically
// applies query timeout to all database operations
func registerQueryTimeoutCallback(db *gorm.DB, dbCfg *config.DBConfig) {
	queryTimeout := time.Duration(dbCfg.QueryTimeoutSecond) * time.Second
	if queryTimeout <= 0 {
		queryTimeout = 30 * time.Second
	}

	applyTimeout := func(db *gorm.DB) {
		if db.Statement != nil && db.Statement.Context != nil {
			if _, hasDeadline := db.Statement.Context.Deadline(); !hasDeadline {
				timeoutCtx, cancel := context.WithTimeout(db.Statement.Context, queryTimeout)
				// Note: GORM doesn't provide a hook to call cancel(),
				// but when the parent context is done, this child is cleaned up.
				_ = cancel
				db.Statement.Context = timeoutCtx
			}
		}
	}

	// Register callback for all query operations
	if err := db.Callback().Query().Before("gorm:query").Register("apply_query_timeout", applyTimeout); err != nil {
		log.Error().Err(err).Msg("failed to register gorm query timeout callback")
	}

	// Register callback for create operations
	if err := db.Callback().Create().Before("gorm:create").Register("apply_query_timeout", applyTimeout); err != nil {
		log.Error().Err(err).Msg("failed to register gorm create timeout callback")
	}

	// Register callback for update operations
	if err := db.Callback().Update().Before("gorm:update").Register("apply_query_timeout", applyTimeout); err != nil {
		log.Error().Err(err).Msg("failed to register gorm update timeout callback")
	}

	// Register callback for delete operations
	if err := db.Callback().Delete().Before("gorm:delete").Register("apply_query_timeout", applyTimeout); err != nil {
		log.Error().Err(err).Msg("failed to register gorm delete timeout callback")
	}
}
