package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vaughan0/go-ini"
	"strings"
	"time"
	"errors"
)

type (
	Config struct {
		MaxIdle     int
		MaxActive   int
		IdleTimeOut time.Duration
		Timeout     time.Duration
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
		Timeout:     time.Second * time.Duration(DefaultTimeout),
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
	if data.Timeout > 0 {
		c.Timeout = data.Timeout
	}
	c.Wait = data.Wait
}

const (
	DefaultNodeName    = "default" // 默认节点配置
	DefaultIdleTimeout = 60        // 连接池默认空余超时时间
	DefaultTimeout     = 60        // dial时的默认请求,读取和写入超时时间
)

var (
	// 节点没有找到
	ErrNodeNotFound = errors.New("node not found")
)

// 读取某个节点的连接池对象,参数node为要读取的节点名称,默认读取default节点
func SelectPool(node ...string) (p *redis.Pool, err error) {
	if len(node) == 0 {
		node = make([]string, 1)
		node[0] = "default"
	}

	p, err = pool.Select(node[0])
	return
}

// 新建连接
func NewConn(nodeName ...string) (conn redis.Conn, err error) {
	p, err := SelectPool(nodeName...)
	if err == nil {
		conn = p.Get()
	}

	return
}

// 执行命令,传入节点名称，执行的命令名，命令参数，返回的第一个参数为返回值，如果出错，第二个参数为空
// 建议用redis.Int()等进行转义结果
func Command(nodeName string, command string, args ...interface{}) (data interface{}, err error) {
	conn, err := NewConn(nodeName)
	defer conn.Close()
	if err != nil {
		return data, err
	}
	data, err = conn.Do(command, args...)
	return
}

// pool 初始化某个node的pool
func (p *Pool) Set(nodeName, nodeConfig string) {
	configSlice := strings.Split(nodeConfig, "?password=")

	p.pools[nodeName] = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: time.Duration(config.IdleTimeOut) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", configSlice[0], config.Timeout*time.Second, config.Timeout*time.Second, config.Timeout*time.Second)
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

// 获取某个连接池
func (p *Pool) Get(nodeName string) (pool *redis.Pool, err error) {
	var exist bool
	if pool, exist = p.pools[nodeName]; !exist {
		return pool, ErrNodeNotFound
	}

	return
}

// 新建连接
func (p *Pool) Conn(nodeName string) (conn redis.Conn, err error) {
	if pool, err := p.Get(nodeName); err == nil {
		conn = pool.Get()
	}

	return
}

// Init 初始化redis配置
func Init(redisConfig Config, nodeConfig ini.File) {
	var (
		findDefaultNode = false
		configNodes     = nodeConfig.Section("redis_nodes")
	)

	config.Set(redisConfig)

	for nodeName, node := range configNodes {
		if nodeName == DefaultNodeName {
			findDefaultNode = true
		}
		pool.Set(nodeName, node)
	}

	if !findDefaultNode {
		pool.Set(DefaultNodeName, "127.0.0.1:6379?password=")
	}
}
