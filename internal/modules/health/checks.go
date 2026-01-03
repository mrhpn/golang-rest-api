package health

import (
	"context"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/app"
)

const healthyStr = "healthy"

func checkDB(ctx context.Context, appCtx *app.Context, checks map[string]string) bool {
	sqlDB, err := appCtx.DB.DB()

	if err != nil {
		checks["db"] = "unhealthy: failed to get database connection"
		return false
	}

	if err = sqlDB.PingContext(ctx); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("database health check failed")
		checks["db"] = "unhealthy: " + err.Error()
		return false
	}

	stats := sqlDB.Stats()
	checks["db"] = healthyStr
	checks["db_open_conns"] = strconv.Itoa(stats.OpenConnections)
	checks["db_idle_conns"] = strconv.Itoa(stats.Idle)
	return true
}

func checkStorage(ctx context.Context, appCtx *app.Context, checks map[string]string) bool {
	if appCtx.MediaService == nil {
		checks["storage"] = "unhealthy: media service not initialized"
		return false
	}

	if err := appCtx.MediaService.HealthCheck(ctx); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("minio health check failed")
		checks["storage"] = "unhealthy: " + err.Error()
		return false
	}

	checks["storage"] = healthyStr
	return true
}

func checkRedis(ctx context.Context, appCtx *app.Context, checks map[string]string) bool {
	if !appCtx.Cfg.Redis.Enabled {
		checks["redis"] = "disabled"
		return true
	}

	if appCtx.Redis == nil {
		checks["redis"] = "unhealthy: redis client not initialized"
		return false
	}

	if err := appCtx.Redis.Ping(ctx).Err(); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("redis health check failed")
		checks["redis"] = "unhealthy: " + err.Error()
		return false
	}

	checks["redis"] = healthyStr
	return true
}
