package redis

import (
	"fmt"

	radix "github.com/mediocregopher/radix/v3"

	"github.com/zoulls/provencal-le-gaulois/config"
)

type Client interface {
	GetDefaultStatus() (*string, error)
}

type client struct {
	config *config.Config
	*radix.Pool
}

func NewClient() (Client, error) {
	config := config.GetConfig()
	pool, err := radix.NewPool("tcp", fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port), int(config.Redis.Pool))
	client := &client{
		config: config,
		Pool:   pool,
	}
	return client, err
}

func (c *client) GetDefaultStatus() (*string, error) {
	var status *string
	err := c.Do(radix.Cmd(status, "GET", "default_status"))
	if err != nil {
		return nil, err
	}
	return status, err
}
