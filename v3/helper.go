package goredis

import (
	"github.com/gomodule/redigo/redis"
)

type Helper struct {
	nodeName string
}

// 新建helper实例
func NewHelper(nodeName ...string) *Helper {
	if len(nodeName) == 0 {
		nodeName = make([]string, 1)
		nodeName[0] = DefaultNodeName
	}
	return &Helper{nodeName: nodeName[0],}
}

// get command
func (h *Helper) Get(key string) (resp interface{}, err error) {
	resp, err = Command(h.nodeName, "GET", key)
	return
}

// set command
func (h *Helper) Set(key string, value interface{}) (err error) {
	_, err = Command(h.nodeName, "SET", key, value)
	return
}

// del command
func (h *Helper) Del(keys ...string) (err error) {
	keysInterface := make([]interface{}, 0, len(keys))
	if len(keys) == 0 {
		return
	}
	for _, v := range keys {
		keysInterface = append(keysInterface, v)
	}
	_, err = Command(h.nodeName, "DEL", keysInterface...)
	return
}

// exists command
func (h *Helper) Exists(key string) (exist bool, err error) {
	exist, err = redis.Bool(Command(h.nodeName, "EXISTS", key))
	return
}

// expire command
func (h *Helper) Expire(key string, seconds int) (err error) {
	_, err = Command(h.nodeName, "EXPIRE", key, seconds)
	return
}

// expireat command
func (h *Helper) ExpireAt(key string, timestamp int64) (err error) {
	_, err = Command(h.nodeName, "EXPIREAT", key, timestamp)
	return
}

// keys command
func (h *Helper) Keys(pattern string) (keys []string, err error) {
	keys, err = redis.Strings(Command(h.nodeName, "KEYS", pattern))
	return
}

// persist command
func (h *Helper) Persist(key string) (err error) {
	if _, err := Command(h.nodeName, "PERSIST", key); err != nil {
		return err
	}
	return nil
}

// ttl command
func (h *Helper) TTL(key string) (ttl int64, err error) {
	ttl, err = redis.Int64(Command(h.nodeName, "TTL", key))
	return
}

// setex command
func (h *Helper) Setex(key string, second int, value interface{}) (err error) {
	_, err = Command(h.nodeName, "SETEX", key, second, value)
	return
}

// setnx command
func (h *Helper) Setnx(key string, value interface{}) (err error) {
	_, err = Command(h.nodeName, "SETNX", key, value)
	return
}

// mset command
func (h *Helper) MSet(valueMap map[string]interface{}) (err error) {
	values := make([]interface{}, 0, len(valueMap)*2)
	for k, v := range valueMap {
		values = append(values, k, v)
	}
	_, err = Command(h.nodeName, "MSET", values...)

	return
}

// mget command
func (h *Helper) MGet(keys ...string) (values []interface{}, err error) {
	if len(keys) == 0 {
		return
	}
	keyInterface := make([]interface{}, 0, len(keys))
	for _, v := range keys {
		keyInterface = append(keyInterface, v)
	}
	values, err = redis.Values(Command(h.nodeName, "MGET", keyInterface...))
	return
}

// decr command
func (h *Helper) Decr(key string) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(h.nodeName, "DECR", key))
	return
}

// decrby command
func (h *Helper) DecrBy(key string, decrNum int) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(h.nodeName, "DECRBY", key, decrNum))
	return
}

// incr command
func (h *Helper) Incr(key string) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(h.nodeName, "INCR", key))
	return
}

// incrby command
func (h *Helper) IncrBy(key string, decrNum int) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(h.nodeName, "INCRBY", key, decrNum))
	return
}

// getset command
func (h *Helper) GetSet(key string, value interface{}) (curValue interface{}, err error) {
	curValue, err = Command(h.nodeName, "GETSET", key, value)
	return
}

// hdel command
func (h *Helper) HDel(key string, fields ...string) (err error) {
	if len(fields) == 0 {
		return
	}
	fieldsInterface := make([]interface{}, 0, len(fields)+1)
	fieldsInterface = append(fieldsInterface, key)
	for _, v := range fields {
		fieldsInterface = append(fieldsInterface, v)
	}
	_, err = Command(h.nodeName, "HDEL", fieldsInterface...)
	return
}

// hexist command
func (h *Helper) HExists(key, field string) (exist bool, err error) {
	exist, err = redis.Bool(Command(h.nodeName, "HEXISTS", key, field))
	return
}

// hset command
func (h *Helper) HSet(key, field string, value interface{}) (err error) {
	_, err = Command(h.nodeName, "HSET", key, field, value)
	return
}

// hmset command
func (h *Helper) HMset(key string, values interface{}) (err error) {
	args := redis.Args{}.Add(key)
	args = args.AddFlat(values)
	_, err = Command(h.nodeName, "HMSET", args...)
	return
}

// hget command
func (h *Helper) HGet(key, field string) (value interface{}, err error) {
	value, err = Command(h.nodeName, "HGET", key, field)
	return
}

// hgetall command
func (h *Helper) HGetAll(key string) (values []interface{}, err error) {
	values, err = redis.Values(Command(h.nodeName, "HGETALL", key))
	return
}

// hkeys command
func (h *Helper) HKeys(key string) (values []string, err error) {
	values, err = redis.Strings(Command(h.nodeName, "HKEYS", key))
	return
}

// hvals command
func (h *Helper) HVals(key string) (values []interface{}, err error) {
	values, err = redis.Values(Command(h.nodeName, "HVALS", key))
	return
}

// hlen command
func (h *Helper) HLen(key string) (length int64, err error) {
	length, err = redis.Int64(Command(h.nodeName, "HLEN", key))
	return
}

// hmget command
func (h *Helper) HMget(key string, fields ...string) (values interface{}, err error) {
	if len(fields) == 0 {
		return
	}
	fieldsInterface := make([]interface{}, 0, len(fields)+1)
	fieldsInterface = append(fieldsInterface, key)
	for _, v := range fields {
		fieldsInterface = append(fieldsInterface, v)
	}
	values, err = Command(h.nodeName, "HMGET", fieldsInterface...)
	return
}
