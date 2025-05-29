package redis

import (
	"fmt"
	envCofig "github.com/rzabhd80/eye-on/internal/envConfig"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	EnvConf *envCofig.AppConfig
}
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func (redisConf *RedisConnection) NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.EnvConf.RedisHost, redisConf.EnvConf.RedisPort),
		Password: redisConf.EnvConf.RedisPassword,
		DB:       redisConf.EnvConf.RedisDB,
	})
}
