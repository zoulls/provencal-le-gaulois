package redis

import (
	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
	"net/url"
	"time"

	radix "github.com/mediocregopher/radix/v3"

	"github.com/zoulls/provencal-le-gaulois/config"
)

type Client interface {
	GetDefaultStatus() (*string, error)
	GetTwitterFollows(follow *config.TwitterFollow) (*string, error)
	Ping() (*string, error)
	Info() (*string, error)
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
	logger.Log.Warn("redis password not configured")
	return radix.Dial(network, addr,
		radix.DialTimeout(30 * time.Second),
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

func (c *client) Ping() (*string, error) {
	var ping string
	err := c.Do(radix.Cmd(&ping, "PING"))
	if err != nil {
		return nil, err
	}
	return &ping, err
}

func (c *client) Info() (*string, error) {
	var infoServ, infoCli, infoKeys string

	err := c.Do(radix.Cmd(&infoServ, "INFO", "server"))
	if err != nil {
		return nil, err
	}
	err = c.Do(radix.Cmd(&infoCli, "INFO", "clients"))
	if err != nil {
		return nil, err
	}
	err = c.Do(radix.Cmd(&infoKeys, "INFO", "Keyspace"))
	if err != nil {
		return nil, err
	}

	info := infoServ+infoCli+infoKeys
	return &info, err
}
