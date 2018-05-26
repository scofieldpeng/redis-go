package goredis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/vaughan0/go-ini"
	"strings"
	"time"
	"errors"
	"math/rand"
	"context"
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
		nodes  map[string]*Node
		config Config
	}
	// 节点
	Node struct {
		pool   *redis.Pool
		slaves []string
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
	pool = Pool{
		nodes: make(map[string]*Node),
	}
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
	// 默认节点配置
	DefaultNodeName = "default"
	// 连接池默认空余超时时间
	DefaultIdleTimeout = 60
	// dial时的默认请求,读取和写入超时时间
	DefaultTimeout = 60

	// node在config中的默认前缀
	NodeConfigPrefix = "redis_node_"
)

var (
	// 节点没有找到
	ErrNodeNotFound = errors.New("node not found")
)

// 获取节点的redis连接池
func (n *Node) GetPool() (p *redis.Pool) {
	p = n.pool
	return
}

// 获取节点的redis连接结构体
// 使用完毕后务必使用conn.Close()释放连接池
func (n *Node) GetConn() (conn redis.Conn) {
	return n.GetPool().Get()
}

// 使用context.Context获取redis的连接
// 使用完毕后务必使用conn.Close()释放连接池
func (n *Node) GetConnContext(ctx context.Context) (conn redis.Conn, err error) {
	return n.GetPool().GetContext(ctx)
}

// 获取当前节点的所有slave节点，如果出错，返回error
func (n *Node) GetSlaves() (nodes []*Node, err error) {
	var node *Node

	nodes = make([]*Node, 0)
	if len(n.slaves) > 0 {
		for _, v := range n.slaves {
			if node, err = GetNode(v); err == nil {
				nodes = append(nodes, node)
			}
		}
	}

	return nodes, err
}

// 获取slave节点,如果不指定slave名称，那么会随机返回一个slave节点，如果查找出错，返回error
// 注意，如果没有找到相关slave节点，将会返回goredis.ErrNodeNotFound错误
func (n *Node) GetSlave(slaveName ...string) (node *Node, err error) {
	if len(n.slaves) == 0 {
		err = ErrNodeNotFound
		return
	}
	if len(slaveName) == 0 {
		slaveName = make([]string, 1)
		// 随机选择一个slave
		rand.Seed(time.Now().UnixNano())
		slaveName[0] = n.slaves[rand.Intn(len(n.slaves))]
	}

	node, err = GetNode(slaveName[0])
	return
}

// 获取节点
func GetNode(node ...string) (n *Node, err error) {
	var exist bool
	if len(node) == 0 {
		node = make([]string, 1)
		node[0] = DefaultNodeName
	}
	n, exist = pool.nodes[node[0]]
	if !exist {
		err = ErrNodeNotFound
	}

	return
}

// 执行命令,传入节点名称，执行的命令名，命令参数，返回的第一个参数为返回值，如果出错，第二个参数为空
// 建议用redis.Int()等进行转义结果
func Command(nodeName string, command string, args ...interface{}) (data interface{}, err error) {
	var n *Node
	n, err = GetNode(nodeName)
	if err != nil {
		return
	}
	conn := n.GetConn()
	defer conn.Close()
	data, err = conn.Do(command, args...)
	return
}

// 在某个node的slave上执行命令，传入节点名称，执行的命令名，命令参数，返回的第一个参数为返回值，如果出错，第二个参数为空
// 建议用redis.Int()等进行转义结果
func CommandOnSlave(nodeName string, command string, args ...interface{}) (data interface{}, err error) {
	var (
		n    *Node
		conn redis.Conn
	)
	n, err = GetNode(nodeName)
	if err != nil {
		return
	}
	n, err = n.GetSlave()
	if err != nil {
		return
	}

	conn = n.GetConn()
	defer conn.Close()

	data, err = conn.Do(command, args...)
	return

}

// pool 初始化某个node的pool
func (p *Pool) SetNode(nodeName, scheme string) {
	if _, exist := p.nodes[nodeName]; !exist {
		p.nodes[nodeName] = &Node{slaves: make([]string, 0)}
	}

	p.nodes[nodeName].pool = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: time.Duration(config.IdleTimeOut) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(scheme,
				redis.DialConnectTimeout(config.IdleTimeOut),
				redis.DialReadTimeout(config.IdleTimeOut),
				redis.DialWriteTimeout(config.IdleTimeOut),
			)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// 设置某个node的slaves，传入slave的node名称即可
func (p *Pool) SetSlaves(nodeName string, slaves []string) {
	if _, exist := p.nodes[nodeName]; !exist {
		p.nodes[nodeName] = &Node{slaves: make([]string, 0)}
	}
	p.nodes[nodeName].slaves = slaves
}

// 获取某个连接池
func (p *Pool) GetNode(nodeName string) (node *Node, err error) {
	var exist bool
	if node, exist = p.nodes[nodeName]; !exist {
		err = ErrNodeNotFound
		return
	}

	return
}

// 新建连接
func (p *Pool) Conn(nodeName string) (conn redis.Conn, err error) {
	var node *Node
	if node, err = p.GetNode(nodeName); err == nil {
		conn = node.GetConn()
	}

	return
}

// Init 初始化redis配置
func Init(redisConfig Config, nodeConfig ini.File) {
	var (
		findDefaultNode = false
		nodeSlaveMap    = make(map[string][]string)
		nodeSchemeMap   = make(map[string]string)

		exist  bool
		scheme string
		slaves string
	)

	config.Set(redisConfig)

	// 遍历出相关的node名称和slave
	for sectionName, section := range nodeConfig {
		if strings.HasPrefix(sectionName, NodeConfigPrefix) {
			nodeName := strings.TrimPrefix(sectionName, NodeConfigPrefix)
			scheme, exist = section["scheme"]
			if exist {
				nodeSchemeMap[nodeName] = scheme
			}
			slaves, exist = section["slaves"]
			if exist {
				nodeSlaveMap[nodeName] = strings.Split(slaves, ",")
			}
		}
	}
	// 过滤掉那些没有dsn的slave的node
	for k, v := range nodeSlaveMap {
		filterSlaves := make([]string, 0, len(v))
		for _, vv := range v {
			if _, exist = nodeSchemeMap[vv]; exist {
				filterSlaves = append(filterSlaves, vv)
			}
		}
		nodeSlaveMap[k] = filterSlaves
	}

	for nodeName, scheme := range nodeSchemeMap {
		if nodeName == DefaultNodeName {
			findDefaultNode = true
		}
		pool.SetNode(nodeName, scheme)
	}
}
