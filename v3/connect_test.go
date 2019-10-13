package goredis

import (
	"github.com/vaughan0/go-ini"
	"testing"
	"github.com/gomodule/redigo/redis"
)

func TestInit(t *testing.T) {
	testIniFile := ini.File{
		"redis_node_default": ini.Section{
			"scheme": "redis://:@localhost:6379",
			"slaves": "slave1",
		},
		"redis_node_slave1": ini.Section{
			"scheme": "redis://:@localhost:6379",
		},
	}

	Init(Config{}, testIniFile)
	node, err := GetNode()
	if err != nil {
		t.Error("get node fail,err:", err.Error())
		return
	}

	conn := node.GetConn()
	defer conn.Close()
	if _, err = conn.Do("SET", "name", "scofield"); err != nil {
		t.Error("set command fail! error: ", err.Error())
		return
	}

	name, err := redis.String(node.Command("GET", "name"))
	if err != nil {
		t.Error("get command fail!error:", err.Error())
	} else if name != "scofield" {
		t.Error("get wrong result,want scofield,get ", name)
	}
	_, err = redis.String(CommandOnSlave(DefaultNodeName, "GET", "name"))
	if err != nil {
		t.Error("get slave command fail! error:", err.Error())
	}

	if _, err = Command(DefaultNodeName, "DEL", "name"); err != nil {
		t.Error("del command fail,error:", err.Error())
		return
	}
	nodes, err := node.GetSlaves()
	if err != nil {
		t.Error("get slaves fai,err:", err.Error())
	} else if len(nodes) != 1 {
		t.Error("get slave count is wrong, want 1 get ", len(nodes))
	}else {
		conn2 := nodes[0].GetConn()
		if _, err := conn.Do("DEL", "name"); err != nil {
			t.Error("run command on slave fail!error:", err.Error())
		}
		defer conn2.Close()
	}
	node, err = node.GetSlave()
	if err != nil {
		t.Error("get slave node fail,err:", err.Error())
		return
	}

	t.Log("active:\t", node.GetPool().ActiveCount())
	t.Log("idle:\t", node.GetPool().IdleCount())
}
