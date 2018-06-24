package goredis

import (
	"testing"
	"github.com/vaughan0/go-ini"
	"github.com/gomodule/redigo/redis"
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
	err := NewHelper().Set("name", "scofield")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGet(t *testing.T) {
	testInit()
	res, err := redis.String(NewHelper().Get("name"))
	if err != nil {
		t.Error(err.Error())
		return
	} else if res != "scofield" {
		t.Error("want scofield, get ", res)
	}
}

func TestDel(t *testing.T) {
	testInit()
	err := NewHelper().Del("name")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestExists(t *testing.T) {
	testInit()
	exist, err := NewHelper().Exists("abc")
	if err != nil {
		t.Error(err.Error())
	} else if exist {
		t.Error("find not exist value")
	}
}

func TestExpire(t *testing.T) {
	testInit()
	err := NewHelper().Expire("name", 1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestExpireAt(t *testing.T) {
	testInit()
	NewHelper().Set("name", "scofield")
	err := NewHelper().ExpireAt("name", time.Now().Unix()+100)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestKeys(t *testing.T) {
	testInit()
	keys, err := NewHelper().Keys("name")
	if err != nil {
		t.Error(err.Error())
	} else if len(keys) > 1 || keys[0] != "name" {
		t.Error("result fail,result:", keys)
	}
}

func TestPersist(t *testing.T) {
	testInit()
	NewHelper().Setex("name", 100, "scofield")
	err := NewHelper().Persist("name")
	if err != nil {
		t.Error(err.Error())
	}
	ttl, err := NewHelper().TTL("name")
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
	ttl, err := NewHelper().TTL("notExistKey")
	if err != nil {
		t.Error(err.Error())
	} else if ttl != -2 {
		t.Error("ttl not exist key not get -2, get: ", ttl)
	}
	NewHelper().Set("name", "scofield")
	ttl, err = NewHelper().TTL("name")
	if err != nil {
		t.Error(err.Error())
	} else if ttl != -1 {
		t.Error("want -1, get ", ttl)
	}
}

func TestSetex(t *testing.T) {
	testInit()
	err := NewHelper().Setex("hello", 1, "world")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSetnx(t *testing.T) {
	testInit()
	err := NewHelper().Setnx("name", "scofield")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMSet(t *testing.T) {
	testInit()
	err := NewHelper().MSet(map[string]interface{}{
		"name": "scofield",
		"age":  26,
	})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMGet(t *testing.T) {
	testInit()
	res, err := NewHelper().MGet("name", "age")
	if err != nil {
		t.Error(err.Error())
	}
	resStr, _ := redis.Strings(res, err)
	t.Log(resStr)
}

func TestDecr(t *testing.T) {
	testInit()
	oldValue, _ := redis.Int64(NewHelper().Get("age"))
	newV, err := NewHelper().Decr("age")
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
	oldValue, _ := redis.Int64(NewHelper().Get("age"))
	newV, err := NewHelper().DecrBy("age", 2)
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
	oldValue, _ := redis.Int64(NewHelper().Get("age"))
	newV, err := NewHelper().Incr("age")
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
	oldValue, _ := redis.Int64(NewHelper().Get("age"))
	newValue, err := NewHelper().IncrBy("age", 2)
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
	oldValue, err := redis.String(NewHelper().GetSet("name", "julia"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("old value:", oldValue)
	newValue, err := redis.String(NewHelper().Get("name"))
	if err != nil {
		t.Error(err.Error())
		return
	} else if newValue != "julia" {
		t.Error("not the new value")
	}
}

func TestHDel(t *testing.T) {
	testInit()
	err := NewHelper().HDel("student", "name", "age")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestHExists(t *testing.T) {
	testInit()
	exist, err := NewHelper().HExists("student", "name")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if exist {
		t.Error("find not exist value")
		return
	}
	NewHelper().HSet("student", "name", "scofield")
	exist, err = NewHelper().HExists("student", "name")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !exist {
		t.Error("can not find exist value")
	}
}

func TestHGet(t *testing.T) {
	testInit()
	err := NewHelper().HSet("student", "name", "scofield")
	if err != nil {
		t.Error(err.Error())
		return
	}
	name, err := redis.String(NewHelper().HGet("student", "name"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if name != "scofield" {
		t.Error("hget wrong value")
	}
}

func TestHGetAll(t *testing.T) {
	type testStudent struct {
		Name string `redis:"name"`
		Age  int    `redis:"age"`
	}
	testInit()
	ts := testStudent{Name: "scofield", Age: 26}
	err := NewHelper().HMset("student", ts)
	if err != nil {
		t.Error(err.Error())
		return
	}
	ts = testStudent{}
	v, err := NewHelper().HGetAll("student")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if err := redis.ScanStruct(v, &ts); err != nil {
		t.Error(err.Error())
		return
	}

	if ts.Name != "scofield" && ts.Age != 26 {
		t.Error("hgetall value is incorrect")
		t.Logf("%#v\n", ts)
	}
}

func TestHKeys(t *testing.T) {
	testInit()
	err := NewHelper().HMset("student", map[string]interface{}{
		"name": "scofield",
		"age":  26,
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	keys, err := NewHelper().HKeys("student")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(keys) != 2 || keys[0] != "name" || keys[1] != "age" {
		t.Error("wrong keys")
		t.Logf("%v", keys)
		return
	}
}

func TestHVals(t *testing.T) {
	testInit()
	NewHelper().HMset("student", map[string]interface{}{
		"name": "scofield",
		"age":  "26",
	})
	values, err := redis.Strings(NewHelper().HVals("student"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(values) != 2 {
		t.Error("wrong values")
		return
	}
	if values[0] != "scofield" && values[1] != "26" {
		t.Error("wrong values")
		return
	}
}

func TestHLen(t *testing.T) {
	testInit()
	l, err := NewHelper().HLen("notExist")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if l != 0 {
		t.Error("should not neq 0")
		return
	}
	if err = NewHelper().HMset("student", map[string]interface{}{
		"name": "scofield",
		"age":  26,
	}); err != nil {
		t.Error(err.Error())
		return
	}

	l, err = NewHelper().HLen("student")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if l != 2 {
		t.Error("invalid hlen,should get 2,get", l)
	}
}

func TestHMget(t *testing.T) {
	testInit()
	err := NewHelper().HMset("student", map[string]interface{}{
		"name": "scofield",
		"age":  26,
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	v, err := redis.Strings(NewHelper().HMget("student", "name"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(v) != 1 {
		t.Error("value length wrong,want 1,get ", len(v))
		return
	}
	if v[0] != "scofield" {
		t.Error("v[0] get wrong,want scofield,get ", v[0])
		return
	}
}
