package cli

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// RedisUserFlag is the ENV var for the Redis username
	RedisUserFlag string = "redis-user"
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
)

// InitRedisFlags initializes RedisFlags command line flags
func InitRedisFlags(flag *pflag.FlagSet) {
	flag.String(RedisUserFlag, "", "Redis username")
	flag.String(RedisPasswordFlag, "", "Redis password")
	flag.String(RedisHostFlag, "localhost", "Redis hostname")
	flag.Int(RedisPortFlag, 6379, "Redis port")
	flag.Int(RedisDBNameFlag, 0, "Redis database")
	flag.Int(RedisConnectTimeoutFlag, 2, "Redis connect timeout in seconds")
}

// CheckRedis validates Redis command line flags
func CheckRedis(v *viper.Viper) error {
	if err := ValidatePort(v, RedisPortFlag); err != nil {
		return err
	}

	environment := v.GetString(EnvironmentFlag)
	if environment == EnvironmentProd {
		if err := ValidateRedisHostInProd(v, RedisHostFlag); err != nil {
			return err
		}
	}

	if environment == EnvironmentTest {
		if err := ValidateRedisHostInTest(v, RedisHostFlag); err != nil {
			return err
		}
	}

	if err := ValidateRedisConnectTimeout(v, RedisConnectTimeoutFlag); err != nil {
		return err
	}

	return nil
}

// InitRedis initializes a Redis pool from command line flags.
// v is the viper Configuration.
// logger is the application logger.
func InitRedis(v *viper.Viper, logger Logger) (*redis.Pool, error) {
	redisUser := v.GetString(RedisUserFlag)
	redisPassword := v.GetString(RedisPasswordFlag)
	redisHost := v.GetString(RedisHostFlag)
	redisPort := strconv.Itoa(v.GetInt(RedisPortFlag))
	redisDBName := v.GetString(RedisDBNameFlag)
	redisConnectTimeout := v.GetInt(RedisConnectTimeoutFlag)
	timeoutDuration := time.Duration(redisConnectTimeout) * time.Second

	s := "redis://%s:%s@%s:%s/%s"
	redisURL := fmt.Sprintf(s, redisUser, redisPassword, redisHost, redisPort, redisDBName)

	if err := testRedisConnection(redisURL, logger, timeoutDuration); err != nil {
		return nil, err
	}

	if err := logRedisConnection(redisURL, logger); err != nil {
		return nil, err
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			connection, err := redis.DialURL(redisURL, redis.DialConnectTimeout(timeoutDuration))
			if err != nil {
				return nil, err
			}
			return connection, nil
		},
	}

	return pool, nil
}

// ValidateRedisHostInProd validates the REDIS_HOST environment variable in prd
func ValidateRedisHostInProd(v *viper.Viper, flagname string) error {
	r := regexp.MustCompile(`([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}`)

	if host := v.GetString(flagname); r.MatchString(host) == false {
		return errors.Errorf("%s can only contain letters, numbers, periods, and dashes", flagname)
	}

	return nil
}

// ValidateRedisHostInTest validates the REDIS_HOST environment variable in test
func ValidateRedisHostInTest(v *viper.Viper, flagname string) error {
	if host := v.GetString(flagname); host != "redis" {
		return errors.New("REDIS_HOST must be 'redis' in the test environment")
	}

	return nil
}

// ValidateRedisConnectTimeout validates a Redis connect timeout passed in from
// the command line
func ValidateRedisConnectTimeout(v *viper.Viper, flagname string) error {
	timeout := v.GetInt(flagname)

	if timeout < 1 || timeout > 5 {
		return errors.Errorf("%s should be between 1 and 5 seconds", flagname)
	}

	return nil
}

func testRedisConnection(redisURL string, logger Logger, timeout time.Duration) error {
	_, err := redis.DialURL(redisURL, redis.DialConnectTimeout(timeout))
	if err != nil {
		errorString := fmt.Sprintf("Failed to connect to Redis after %s", timeout)
		logger.Error(errorString, zap.Error(err))
		return err
	}

	return nil
}

func logRedisConnection(redisURL string, logger Logger) error {
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		return errors.Errorf("%s is an invalid redis URL", redisURL)
	}
	password, _ := parsedURL.User.Password()
	var maskedURL string
	if password != "" {
		maskedURL = strings.Replace(redisURL, password, "*****", 1)
	} else {
		maskedURL = redisURL
	}
	logger.Info("Connecting to Redis", zap.String("url", maskedURL))

	return nil
}
