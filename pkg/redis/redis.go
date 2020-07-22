package redis

import (
	"fmt"
	"time"

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
	pool, err := radix.NewPool("tcp", fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port), int(conf.Redis.Pool), radix.PoolConnFunc(customConnFunc))
	client := &client{
		config: conf,
		Pool:   pool,
	}
	return client, err
}

// Custom function with Auth connection
func customConnFunc(network, addr string) (radix.Conn, error) {
	conf := config.GetConfig()
	return radix.Dial(network, addr,
		radix.DialTimeout(30 * time.Second),
		radix.DialAuthPass(conf.Redis.Pass),
	)
}

func (c *client) GetDefaultStatus() (*string, error) {
	var status string
	err := c.Do(radix.Cmd(&status, "GET", "defaultStatus"))
	if err != nil {
		return nil, err
	}
	return &status, err
}

func (c *client) GetTwitterFollows(follow *config.TwitterFollow) (*string, error) {
	var list string
	err := c.Do(radix.Cmd(&list, "GET", follow.Key))
	if err != nil {
		return nil, err
	}
	return &list, err
}