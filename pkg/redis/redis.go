package redis

import (
	"github.com/mediocregopher/radix"

	"github.com/zoulls/provencal-le-gaulois/config"
)

func NewRedisClient(conf config.RedisConfig) error {
	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 10)
	if err != nil {
		return err
	}
}
