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
			"slave":  "slave1",
		},
		"redis_node_slave1": ini.Section{
			"scheme": "redis://:@localhost:6379",
		},
	}

	Init(Config{}, testIniFile)
	node, err := GetNode("default")
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

	if name, err := redis.String(Command("default", "GET", "name")); err != nil {
		t.Error("get command fail, error:", err.Error())
		return
	} else if name != "scofield" {
		t.Error("get command resutl fail, result:", name, ", want scofield")
	}
}
