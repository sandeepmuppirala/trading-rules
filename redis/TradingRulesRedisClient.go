package redis

import (
	log "github.com/sirupsen/logrus"
	"os"
	"trading-rules/constants"

	"github.com/gomodule/redigo/redis"
)

var connPool *redis.Pool

// InitRedisConnection : Initiate Redis connection pool
func InitRedisConnection() {
	connPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(constants.RedisDialProtocol, constants.RedisHostname, 
													redis.DialPassword(constants.RedisPassword),
															redis.DialUsername(constants.RedisUsername))
			if err != nil {
				log.Debug("ERROR: failed to initiate redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

// GetOperation : Redis GET operation
func GetOperation(key string) (string, error) {
	conn := connPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return s, nil
}

// SetOperation : Redis SET operation
func SetOperation(key string, val string) error {
	conn := connPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		return err
	}

	return nil
}

// DelOperation : Redis DEL operation
func DelOperation(key string, val string) error {
	conn := connPool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key, "")
	if err != nil {
		return err
	}

	return nil
}