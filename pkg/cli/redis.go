package cli

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// RedisPasswordFlag is the ENV var for the Redis password
	RedisPasswordFlag string = "redis-password"
	// RedisHostFlag is the ENV var for the Redis hostname
	RedisHostFlag string = "redis-host"
	// RedisPortFlag is the ENV var for the Redis port
	RedisPortFlag string = "redis-port"
	// RedisDBNameFlag is the ENV var for the Redis database name, which
	// is represented by a positive integer. Using multiple databases in
	// the same Redis instance allows separating concerns.
	RedisDBNameFlag string = "redis-db-name"
	// RedisConnectTimeoutFlag specifies how long to wait to establish a
	// connection to the Redis instance
	RedisConnectTimeoutFlag string = "redis-connect-timeout-in-seconds"
	// RedisEnabledFlag specifies whether or not we attempt to connect
	// to Redis. For example, apps that use mTLS don't need Redis.
	RedisEnabledFlag string = "redis-enabled"
	// RedisSSLEnabledFlag specifies if SSL mode is enabled for connections
	RedisSSLEnabledFlag string = "redis-ssl-enabled"
	// RedisMaxIdleFlag specifies the maximum number of idle connections in the pool
	RedisMaxIdleFlag string = "redis-max-idle"
	// RedisIdleTimeoutFlag Closes connections after this duration
	RedisIdleTimeoutFlag string = "redis-idle-timeout"
)

// InitRedisFlags initializes RedisFlags command line flags
func InitRedisFlags(flag *pflag.FlagSet) {
	flag.String(RedisPasswordFlag, "", "Redis password")
	flag.String(RedisHostFlag, "localhost", "Redis hostname")
	flag.Int(RedisPortFlag, 6379, "Redis port")
	flag.Int(RedisDBNameFlag, 0, "Redis database")
	flag.Duration(RedisConnectTimeoutFlag, 2*time.Second, "Redis connect timeout in seconds")
	flag.Int(RedisMaxIdleFlag, 10, "Redis maximum number of idle connections in the pool")
	flag.Duration(RedisIdleTimeoutFlag, 240*time.Second, "Redis idle timeout in seconds")
	flag.Bool(RedisEnabledFlag, true, "Whether or not Redis is enabled")
	flag.Bool(RedisSSLEnabledFlag, false, "Whether or not Redis SSL is enabled")
}

// CheckRedis validates Redis command line flags
func CheckRedis(v *viper.Viper) error {
	enabled := v.GetBool(RedisEnabledFlag)
	if !enabled {
		return nil
	}

	if err := ValidatePort(v, RedisPortFlag); err != nil {
		return err
	}

	if err := ValidateHost(v, RedisHostFlag); err != nil {
		return err
	}

	connectTimeout := v.GetDuration(RedisConnectTimeoutFlag)
	if connectTimeout < 1*time.Second || connectTimeout > 5*time.Second {
		return errors.Errorf("%s should be between 1 and 5 seconds", RedisConnectTimeoutFlag)
	}

	maxIdle := v.GetInt(RedisMaxIdleFlag)
	if maxIdle < 1 || maxIdle > 20 {
		return errors.Errorf("%s should be between 1 and 20", RedisMaxIdleFlag)
	}

	idleTimeout := v.GetDuration(RedisIdleTimeoutFlag)
	if idleTimeout < 30*time.Second || idleTimeout > 300*time.Second {
		return errors.Errorf("%s should be between 30 and 300 seconds", RedisIdleTimeoutFlag)
	}

	return nil
}

// InitRedis initializes a Redis pool from command line flags.
// v is the viper Configuration.
// logger is the application logger.
func InitRedis(v *viper.Viper, logger Logger) (*redis.Pool, error) {
	enabled := v.GetBool(RedisEnabledFlag)
	if !enabled {
		return nil, nil
	}

	redisPassword := v.GetString(RedisPasswordFlag)
	redisHost := v.GetString(RedisHostFlag)
	redisPort := v.GetInt(RedisPortFlag)
	redisDBName := v.GetInt(RedisDBNameFlag)
	redisConnectTimeout := v.GetDuration(RedisConnectTimeoutFlag)
	redisSSLEnabled := v.GetBool(RedisSSLEnabledFlag)
	redisMaxIdle := v.GetInt(RedisMaxIdleFlag)
	redisIdleTimeout := v.GetDuration(RedisIdleTimeoutFlag)

	// Log the redis URI
	s := "redis://:%s@%s:%d?db=%d"
	redisURI := fmt.Sprintf(s, "*****", redisHost, redisPort, redisDBName)
	if redisPassword == "" {
		s = "redis://%s:%d?db=%d"
		redisURI = fmt.Sprintf(s, redisHost, redisPort, redisDBName)
	}
	logger.Info("Connecting to Redis", zap.String("url", redisURI))

	// Configure Redis TLS Config
	redisTLSConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Redis Dial requires a minimal URI containing just the host and port
	redisURLTemplate := "%s:%s"
	redisURL := fmt.Sprintf(redisURLTemplate, redisHost, strconv.Itoa(redisPort))

	if testRedisErr := testRedisConnection(redisURL, redisPassword, redisDBName, redisConnectTimeout, redisSSLEnabled, &redisTLSConfig, logger); testRedisErr != nil {
		return nil, testRedisErr
	}

	pool := &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeout,
		Dial: func() (redis.Conn, error) {
			connection, connectionErr := redis.Dial(
				"tcp",
				redisURL,
				redis.DialDatabase(redisDBName),
				redis.DialPassword(redisPassword),
				redis.DialConnectTimeout(redisConnectTimeout),
				redis.DialUseTLS(redisSSLEnabled),
				redis.DialTLSConfig(&redisTLSConfig),
			)
			if connectionErr != nil {
				return nil, connectionErr
			}
			return connection, nil
		},
	}

	return pool, nil
}

func testRedisConnection(redisURL, redisPassword string, redisDBName int, redisConnectTimeout time.Duration, redisSSLEnabled bool, redisTLSConfig *tls.Config, logger Logger) error {
	// Confirm the connection works
	logger.Info("Testing Redis connection...")

	redisConnection, redisConnectionErr := redis.Dial(
		"tcp",
		redisURL,
		redis.DialDatabase(redisDBName),
		redis.DialPassword(redisPassword),
		redis.DialConnectTimeout(redisConnectTimeout),
		redis.DialUseTLS(redisSSLEnabled),
		redis.DialTLSConfig(redisTLSConfig),
	)
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to close an asynchronous connection
	//RA: Given the functions causing the lint errors are used close an asynchronous connection in order to prevent it
	//RA: from running indefinitely, it is not deemed a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	defer redisConnection.Close() // nolint:errcheck

	errorString := fmt.Sprintf("Failed to connect to Redis after %s", redisConnectTimeout)
	var finalErrorString string
	if redisPassword == "" {
		finalErrorString = errorString + ". No password provided."
	}

	if redisConnectionErr != nil {
		logger.Error(finalErrorString, zap.Error(redisConnectionErr))
		return redisConnectionErr
	}
	logger.Info("...Redis connection successful!")

	logger.Info("Starting Redis ping...")
	_, pingErr := redis.String(redisConnection.Do("PING"))
	if pingErr != nil {
		logger.Error("failed to ping Redis", zap.Error(pingErr))
		return pingErr
	}

	logger.Info("...Redis ping successful!")

	return nil
}
