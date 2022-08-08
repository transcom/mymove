package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

func redisHealthCheck(pool *redis.Pool, logger *zap.Logger, data map[string]interface{}) map[string]interface{} {
	conn := pool.Get()

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			logger.Error("Failed to close redis connection", zap.Error(closeErr))
		}
	}()

	pong, err := redis.String(conn.Do("PING"))
	if err != nil {
		logger.Error("Failed to ping Redis during health check", zap.Error(err))
	}
	logger.Info("Health check Redis ping", zap.String("ping_response", pong))

	data["redis"] = err == nil

	return data
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
			}
			data["database"] = dbErr == nil
			if redisPool != nil {
				data = redisHealthCheck(redisPool, appCtx.Logger(), data)
			}
		}
		newEncoderErr := json.NewEncoder(w).Encode(data)
		if newEncoderErr != nil {
			appCtx.Logger().Error("Failed encoding health check response", zap.Error(newEncoderErr))
		}

		// We are not using request middleware here so logging directly in the check
		var protocol string
		if r.TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}

		fields := []zap.Field{
			zap.String("accepted-language", r.Header.Get("accepted-language")),
			zap.Int64("content-length", r.ContentLength),
			zap.String("host", r.Host),
			zap.String("method", r.Method),
			zap.String("protocol", protocol),
			zap.String("protocol-version", r.Proto),
			zap.String("referer", r.Header.Get("referer")),
			zap.String("source", r.RemoteAddr),
			zap.String("url", r.URL.String()),
			zap.String("user-agent", r.UserAgent()),
		}

		// Append x- headers, e.g., x-forwarded-for.
		for name, values := range r.Header {
			if nameLowerCase := strings.ToLower(name); strings.HasPrefix(nameLowerCase, "x-") {
				if len(values) > 0 {
					fields = append(fields, zap.String(nameLowerCase, values[0]))
				}
			}
		}

		// Log the number of headers, which can be used for finding abnormal requests
		fields = append(fields, zap.Int("headers", len(r.Header)))

		appCtx.Logger().Info("Request health", fields...)

	}
}
