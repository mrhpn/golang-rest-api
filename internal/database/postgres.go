package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

const (
	maxElapsedTimeMillisecond = 200 * time.Millisecond
	timeoutSecond             = 5 * time.Second
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

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		// Log errors with full context
		lgr.Error().
			Err(err).
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Msg(sql)
	case elapsed > maxElapsedTimeMillisecond:
		// Log slow queries as warnings
		lgr.Warn().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("perf", "SLOW_QUERY").
			Msgf("%s", sql)
	default:
		// Log normal queries as debug
		lgr.Debug().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Msg(sql)
	}
}

// Connect establishes a database connection with retry logic and configurable pool settings
func Connect(parentCtx context.Context, dsn string, dbCfg *config.DBConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		AllowGlobalUpdate: false, // safety: prevent global updates/deletes without a WHERE clause
		PrepareStmt:       true,  // performance: cache prepared statements
	}

	gormConfig.Logger = &CustomGormLogger{}

	db, err := connectWithRetry(dsn, gormConfig, dbCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	poolCfg := normalizePoolConfig(dbCfg)
	applyPoolConfig(sqlDB, poolCfg)

	// Verify connection
	pingCtx, cancel := context.WithTimeout(context.Background(), timeoutSecond)
	defer cancel()

	if err = sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Int("max_open_conns", poolCfg.maxOpenConns).
		Int("max_idle_conns", poolCfg.maxIdleConns).
		Dur("conn_max_lifetime", poolCfg.connMaxLifetime).
		Dur("conn_max_idle_time", poolCfg.connMaxIdleTime).
		Msg("✅ Database — connected successfully")

	// Register callback to automatically apply query timeout to all queries
	registerQueryTimeoutCallback(db, dbCfg)

	// Start periodic connection pool metrics logging
	if dbCfg.DBPoolMetricsEnabled {
		metricsCtx := parentCtx
		if metricsCtx == nil {
			metricsCtx = context.Background()
		}
		interval := time.Duration(dbCfg.DBPoolMetricsLogIntervalSecond) * time.Second
		if interval <= 0 {
			interval = time.Duration(constants.DBPoolMetricsLogIntervalSecond) * time.Second
		}
		startConnectionPoolMetrics(metricsCtx, sqlDB, interval)
	}

	return db, nil
}

func connectWithRetry(dsn string, gormConfig *gorm.Config, dbCfg *config.DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxAttempts := dbCfg.RetryAttempts
	if maxAttempts <= 0 {
		maxAttempts = constants.DBMaxRetryAttempts
	}

	retryDelay := time.Duration(dbCfg.RetryDelaySecond) * time.Second
	if retryDelay <= 0 {
		retryDelay = constants.DBRetryDelaySecond * time.Second
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

	return db, nil
}

type poolConfig struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

func normalizePoolConfig(dbCfg *config.DBConfig) poolConfig {
	maxOpenConns := dbCfg.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = constants.DBMaxOpenConns
	}

	maxIdleConns := dbCfg.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = constants.DBMaxIdleConns
	}

	connMaxLifetime := time.Duration(dbCfg.ConnMaxLifetimeMinute) * time.Minute
	if connMaxLifetime <= 0 {
		connMaxLifetime = constants.DBMaxLifetimeMinute * time.Minute
	}

	connMaxIdleTime := time.Duration(dbCfg.ConnMaxIdleTimeMinute) * time.Minute
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = constants.DBConnMaxIdleTimeMinute * time.Minute
	}

	return poolConfig{
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
		connMaxIdleTime: connMaxIdleTime,
	}
}

func applyPoolConfig(sqlDB *sql.DB, cfg poolConfig) {
	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.connMaxIdleTime)
}

// registerQueryTimeoutCallback registers a GORM callback that automatically
// applies query timeout to all database operations
func registerQueryTimeoutCallback(db *gorm.DB, dbCfg *config.DBConfig) {
	queryTimeout := time.Duration(dbCfg.QueryTimeoutSecond) * time.Second
	if queryTimeout <= 0 {
		queryTimeout = constants.DBMaxQueryTimeoutSecond * time.Second
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

// startConnectionPoolMetrics starts a goroutine that periodically logs database
// connection pool statistics for monitoring and debugging purposes.
func startConnectionPoolMetrics(ctx context.Context, sqlDB *sql.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		log.Info().Dur("interval", interval).Msg("📝 Database — pool metrics logging started")
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("✓ Database pool metrics logger stopped")
				return
			case <-ticker.C:
				stats := sqlDB.Stats()
				log.Info().
					Int("open_connections", stats.OpenConnections).
					Int("in_use", stats.InUse).
					Int("idle", stats.Idle).
					Int64("wait_count", stats.WaitCount).
					Dur("wait_duration", stats.WaitDuration).
					Msg("Database connection pool stats")
			}
		}
	}()
}
