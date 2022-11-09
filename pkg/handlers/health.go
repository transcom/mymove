package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

func redisHealthCheck(pool *redis.Pool, logger *zap.Logger) error {
	conn := pool.Get()

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			logger.Error("Failed to close redis connection", zap.Error(closeErr))
		}
	}()

	pong, err := redis.String(conn.Do("PING"))
	if err != nil {
		logger.Error("Failed to ping Redis during health check", zap.Error(err))
	} else {
		logger.Info("Health check Redis ping", zap.String("ping_response", pong))
	}

	return err
}

func healthCheckError(appCtx appcontext.AppContext, w http.ResponseWriter, data map[string]interface{}) {
	healthCheck, err := json.Marshal(data)
	if err != nil {
		appCtx.Logger().Warn("Cannot marshal health data", zap.Error(err))
		healthCheck = []byte("{}")
	}
	http.Error(w, string(healthCheck), http.StatusInternalServerError)
}

// NewHealthHandler creates a http.HandlerFunc for a health endpoint.
// If redisPool is nil, redis health will not be checked
func NewHealthHandler(appCtx appcontext.AppContext, redisPool *redis.Pool, gitBranch string, gitCommit string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"gitBranch": gitBranch,
			"gitCommit": gitCommit,
		}
		// Check and see if we should disable DB query with '?database=false'
		// Disabling the DB is useful for Route53 health checks which require the TLS
		// handshake be less than 4 seconds and the status code return in less than
		// two seconds. https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/dns-failover-determining-health-of-endpoints.html
		showDB, ok := r.URL.Query()["database"]

		// Always show DB unless key set to "false"
		if !ok || (ok && showDB[0] != "false") {
			appCtx.Logger().Info("Health check connecting to the DB")
			dbErr := appCtx.DB().RawQuery("SELECT 1;").Exec()
			if dbErr != nil {
				appCtx.Logger().Error("Failed database health check", zap.Error(dbErr))
				data["database"] = false
				healthCheckError(appCtx, w, data)
				return
			}
			data["database"] = true
			if redisPool != nil {
				redisErr := redisHealthCheck(redisPool, appCtx.Logger())
				if redisErr != nil {
					data["redis"] = false
					healthCheckError(appCtx, w, data)
					return
				}
				data["redis"] = true
			}
		}
		newEncoderErr := json.NewEncoder(w).Encode(data)
		if newEncoderErr != nil {
			appCtx.Logger().Error("Failed encoding health check response", zap.Error(newEncoderErr))
			http.Error(w, "failed health check", http.StatusInternalServerError)
			return
		}

		appCtx.Logger().Info("Request health ok")
	}
}
