package redis

import (
	"git.name.im/bamboo/config.git"
	"github.com/garyburd/redigo/redis"
	"github.com/vaughan0/go-ini"
	"strings"
	"time"
)

var (
	configCaches ini.File               // 配置缓存
	pools        map[string]*redis.Pool // redis进程池
)

const(
    DefaultNodeName = "default" // 默认节点配置
	DefaultIdleTimeout = 60     // 连接池默认空余超时时间
	DefaultTimeout = 60         // dial时的默认请求,读取和写入超时时间
)

// 读取某个节点的内存池对象,参数node为要读取的节点名称,默认读取default节点
func Pool(node ...string) *redis.Pool {
	if len(node) == 0 {
		node = make([]string, 1)
		node[0] = "default"
	}
	return pools[node[0]]
}

// pool 初始化某个node的pool
func pool(nodeName, nodeConfig string) {
	idleNum := config.Int(configCaches.Get("config", "maxIdle"))
	if idleNum < 1 {
		idleNum = 5
	}

	idleTimeout := config.Int(configCaches.Get("config", "idleTimeout"))
	if idleTimeout < 1 {
		idleTimeout = DefaultIdleTimeout
	}
	timeout := config.Int(configCaches.Get("config","timeout"))
	if timeout < 1 {
		timeout = DefaultTimeout
	}

	configSlice := strings.Split(nodeConfig, "?password=")

	pools[nodeName] = &redis.Pool{
		MaxIdle:     idleNum,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", configSlice[0],time.Duration(timeout)*time.Second,time.Duration(timeout)*time.Second,time.Duration(timeout)*time.Second)
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

// Init 初始化redis配置
func Init(config ini.File) {
	pools = make(map[string]*redis.Pool)

	configCaches = config
    configNodes := configCaches.Section("nodes")
    findDefaultNode := false
    for nodeName,node := range configNodes {
         if nodeName == DefaultNodeName {
			 findDefaultNode = true
		 }

		pool(nodeName,node)
    }

	if !findDefaultNode {
		pool(DefaultNodeName,"127.0.0.1:6379?password=")
	}
}
