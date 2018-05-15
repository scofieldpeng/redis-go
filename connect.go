package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vaughan0/go-ini"
	"strings"
	"time"
	configService "github.com/scofieldpeng/config-go"
)

type (
	Config struct {
		MaxIdle     int
		MaxActive   int
		IdleTimeOut time.Duration
		Wait        bool
	}
	// 连接池结构体
	Pool struct {
		pools  map[string]*redis.Pool
		config Config
	}
)

var (
	config = Config{
		MaxIdle:     5,
		MaxActive:   100,
		IdleTimeOut: time.Second * time.Duration(DefaultIdleTimeout),
		Wait:        true,
	}
	pool Pool
)

func (c *Config) Set(data Config) {
	if data.MaxIdle > 0 {
		c.MaxIdle = data.MaxIdle
	}
	if data.MaxActive > 0 {
		c.MaxActive = data.MaxActive
	}
	if data.IdleTimeOut > 0 {
		c.IdleTimeOut = data.IdleTimeOut
	}
	c.Wait = data.Wait
}

const (
	DefaultNodeName    = "default" // 默认节点配置
	DefaultIdleTimeout = 60        // 连接池默认空余超时时间
	DefaultTimeout     = 60        // dial时的默认请求,读取和写入超时时间
)

// 读取某个节点的内存池对象,参数node为要读取的节点名称,默认读取default节点
func Select(node ...string) *redis.Pool {
	if len(node) == 0 {
		node = make([]string, 1)
		node[0] = "default"
	}
	return pools[node[0]]
}

// pool 初始化某个node的pool
func (p *Pool) pool(nodeName, nodeConfig string) {
	idleNum := configService.Int(configCaches.Get("config", "maxIdle"))
	if idleNum < 1 {
		idleNum = 5
	}

	idleTimeout := config.Int(configCaches.Get("config", "idleTimeout"))
	if idleTimeout < 1 {
		idleTimeout = DefaultIdleTimeout
	}
	timeout := config.Int(configCaches.Get("config", "timeout"))
	if timeout < 1 {
		timeout = DefaultTimeout
	}

	configSlice := strings.Split(nodeConfig, "?password=")

	pools[nodeName] = &redis.Pool{
		MaxIdle:     idleNum,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", configSlice[0], time.Duration(timeout)*time.Second, time.Duration(timeout)*time.Second, time.Duration(timeout)*time.Second)
			if err != nil {
				return nil, err
			}
			if len(configSlice) == 2 && configSlice[1] != "" {
				if _, err := c.Do("AUTH", configSlice[1]); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (p *Pool) Select(nodeName string)(pool *redis.Pool,err error) {
	return
}

// Init 初始化redis配置
func Init(redisConfig Config, nodeConfig ini.File) {
	var (
		findDefaultNode = false
		configNodes     = nodeConfig.Section("redis_nodes")
	)

	config.Set(redisConfig)
	pools = make(map[string]*redis.Pool)

	for nodeName, node := range configNodes {
		if nodeName == DefaultNodeName {
			findDefaultNode = true
		}

		pool(nodeName, node)
	}

	if !findDefaultNode {
		pool(DefaultNodeName, "127.0.0.1:6379?password=")
	}
}
