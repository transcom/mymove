package cli

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// RedisURLFlag is the ENV var for the Redis URL
	RedisURLFlag string = "redis-url"
	// RedisPortFlag is the ENV var for the Redis port
	RedisPortFlag string = "redis-port"
	// RedisDBFlag is the ENV var for the Redis database
	RedisDBFlag string = "redis-db"
)

// InitRedisFlags initializes RedisFlags command line flags
func InitRedisFlags(flag *pflag.FlagSet) {
	flag.String(RedisURLFlag, "redis://localhost", "Redis URL with username, password, and hostname")
	flag.String(RedisPortFlag, "6379", "Redis port")
	flag.String(RedisDBFlag, "0", "Redis database")
}

// CheckRedis validates Redis command line flags
func CheckRedis(v *viper.Viper) error {
	if v.GetString(RedisURLFlag) == "redis://localhost" {
		return nil
	}

	r, _ := regexp.Compile(`redis://\w+:\w+@\w[-.\w]+:\d`)

	if r.MatchString(RedisURLFlag) == false {
		return errors.Errorf("%s must follow the scheme 'redis://username:password@hostname'", RedisURLFlag)
	}

	return nil
}
