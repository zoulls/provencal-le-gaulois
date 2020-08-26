package redis

import (
	"errors"
	"net/url"
	"time"

	radix "github.com/mediocregopher/radix/v3"

	"github.com/zoulls/provencal-le-gaulois/config"
)

type Client interface {
	GetDefaultStatus() (*string, error)
	GetTwitterFollows(follow *config.TwitterFollow) (*string, error)
}

type client struct {
	config *config.Config
	*radix.Pool
}

func NewClient() (Client, error) {
	conf := config.GetConfig()
	pool, err := radix.NewPool("tcp", conf.Redis.URL, int(conf.Redis.Pool), radix.PoolConnFunc(customConnFunc))
	client := &client{
		config: conf,
		Pool:   pool,
	}
	return client, err
}

// Custom function with Auth connection
func customConnFunc(network, addr string) (radix.Conn, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	pass, exists := u.User.Password()
	if exists {
		return radix.Dial(network, addr,
			radix.DialTimeout(30 * time.Second),
			radix.DialAuthPass(pass),
		)
	}
	return nil, errors.New("redis password not configured")
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
