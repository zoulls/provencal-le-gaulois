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
	conf := config.GetConfig()
	// TODO add user/pass
	pool, err := radix.NewPool("tcp", fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port), int(conf.Redis.Pool))
	client := &client{
		config: conf,
		Pool:   pool,
	}
	return client, err
}

func (c *client) GetDefaultStatus() (*string, error) {
	var status string
	err := c.Do(radix.Cmd(&status, "GET", "default_status"))
	if err != nil {
		return nil, err
	}
	return &status, err
}
