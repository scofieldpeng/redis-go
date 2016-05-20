package redis

import (
	"github.com/vaughan0/go-ini"
	"testing"
    "github.com/garyburd/redigo/redis"
)

func TestInit(t *testing.T) {
	testIniFile := ini.File{
		"config": ini.Section{
			"maxIdle":     "5",
			"idleTimeout": "10",
			"timeout":"5",
		},
		"nodes": ini.Section{
			"default2": "127.2.1.1:6379",
		},
	}

	Init(testIniFile)

	conn := Pool().Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("SET", "test", 1))
	if err != nil {
		t.Error("exec command `set test 1` fail,error:", err.Error())
	} else if res != "OK" {
		t.Error("exec command `set test 1` fail,want to get the result `OK`,but get `", res,"`")
	}

	resInt, err := redis.Int(conn.Do("GET", "test"))
	if err != nil {
		t.Error("exec command `get test` fail,error:", err.Error())
	} else if resInt != 1 {
        t.Error("exec command `get test` fail,want to get the result 1,but get ",resInt)
    }

    conn2 := Pool("default2").Get()
    defer conn2.Close()

    if _,err := conn2.Do("ping");err == nil {
        t.Error("conect to a not exist redis-server should fail!")
    }
}
