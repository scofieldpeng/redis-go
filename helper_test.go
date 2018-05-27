package goredis

import (
	"testing"
	"github.com/vaughan0/go-ini"
	"github.com/garyburd/redigo/redis"
	"time"
)

// init just 4 test
func testInit() {
	Init(Config{}, ini.File{
		"redis_node_default": ini.Section{
			"scheme": "redis://:@localhost:6379",
		},
	})
}

func TestSet(t *testing.T) {
	testInit()
	err := Set("default", "name", "scofield")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGet(t *testing.T) {
	testInit()
	res, err := redis.String(Get("default", "name"))
	if err != nil {
		t.Error(err.Error())
		return
	} else if res != "scofield" {
		t.Error("want scofield, get ", res)
	}
}

func TestDel(t *testing.T) {
	testInit()
	err := Del("default", "name")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestExists(t *testing.T) {
	testInit()
	exist, err := Exists("default", "abc")
	if err != nil {
		t.Error(err.Error())
	} else if exist {
		t.Error("find not exist value")
	}
}

func TestExpire(t *testing.T) {
	testInit()
	err := Expire("default", "name", 1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestExpireAt(t *testing.T) {
	testInit()
	Set("default", "name", "scofield")
	err := ExpireAt("default", "name", time.Now().Unix()+100)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestKeys(t *testing.T) {
	testInit()
	keys, err := Keys(DefaultNodeName, "name")
	if err != nil {
		t.Error(err.Error())
	} else if len(keys) > 1 || keys[0] != "name" {
		t.Error("result fail,result:", keys)
	}
}

func TestPersist(t *testing.T) {
	testInit()
	Setex(DefaultNodeName, "name", 100, "scofield")
	err := Persist(DefaultNodeName, "name")
	if err != nil {
		t.Error(err.Error())
	}
	ttl, err := TTL(DefaultNodeName, "name")
	if err != nil {
		t.Error("get ttl fail,err:", err.Error())
		return
	}
	if ttl != -1 {
		t.Error("want ttl is -1,get ", ttl)
	}
}

func TestTTL(t *testing.T) {
	testInit()
	ttl, err := TTL(DefaultNodeName, "notExistKey")
	if err != nil {
		t.Error(err.Error())
	} else if ttl != -2 {
		t.Error("ttl not exist key not get -2, get: ", ttl)
	}
	Set(DefaultNodeName, "name", "scofield")
	ttl, err = TTL(DefaultNodeName, "name")
	if err != nil {
		t.Error(err.Error())
	} else if ttl != -1 {
		t.Error("want -1, get ", ttl)
	}
}

func TestSetex(t *testing.T) {
	testInit()
	err := Setex(DefaultNodeName, "hello", 1, "world")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSetnx(t *testing.T) {
	testInit()
	err := Setnx(DefaultNodeName, "name", "scofield")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMSet(t *testing.T) {
	testInit()
	err := MSet(DefaultNodeName, map[string]interface{}{
		"name": "scofield",
		"age":  26,
	})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMGet(t *testing.T) {
	testInit()
	res, err := MGet(DefaultNodeName, "name", "age")
	if err != nil {
		t.Error(err.Error())
	}
	resStr, _ := redis.Strings(res, err)
	t.Log(resStr)
}

func TestDecr(t *testing.T) {
	testInit()
	oldValue, _ := redis.Int64(Get(DefaultNodeName, "age"))
	newV, err := Decr(DefaultNodeName, "age")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if newV != oldValue-1 {
		t.Error("new value not decr value")
	}
}

func TestDecrBy(t *testing.T) {
	testInit()
	oldValue, _ := redis.Int64(Get(DefaultNodeName, "age"))
	newV, err := DecrBy(DefaultNodeName, "age", 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if newV != oldValue-2 {
		t.Error("new value not decr value")
	}
}

func TestIncr(t *testing.T) {
	testInit()
	oldValue, _ := redis.Int64(Get(DefaultNodeName, "age"))
	newV, err := Incr(DefaultNodeName, "age")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if newV != oldValue+1 {
		t.Error("new vlaue not incr value")
	}
}

func TestIncrBy(t *testing.T) {
	testInit()
	oldValue, _ := redis.Int64(Get(DefaultNodeName, "age"))
	newValue, err := IncrBy(DefaultNodeName, "age", 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if newValue != oldValue+2 {
		t.Error("new value not incr value")
	}
}

func TestGetSet(t *testing.T) {
	testInit()
	oldValue, err := redis.String(GetSet(DefaultNodeName, "name", "julia"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("old value:", oldValue)
	newValue, err := redis.String(Get(DefaultNodeName, "name"))
	if err != nil {
		t.Error(err.Error())
		return
	} else if newValue != "julia" {
		t.Error("not the new value")
	}
}

func TestHDel(t *testing.T) {
	testInit()
	err := HDel(DefaultNodeName, "student", "name", "age")
	if err != nil {
		t.Error(err.Error())
	}
}