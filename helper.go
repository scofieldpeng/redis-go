package goredis

import (
	"github.com/garyburd/redigo/redis"
)

// get command
func Get(nodeName, key string) (resp interface{}, err error) {
	resp, err = Command(nodeName, "GET", key)
	return
}

// set command
func Set(nodeName, key string, value interface{}) (err error) {
	_, err = Command(nodeName, "SET", key, value)
	return
}

// del command
func Del(nodeName string, keys ...string) (err error) {
	keysInterface := make([]interface{}, 0, len(keys))
	if len(keys) == 0 {
		return
	}
	for _, v := range keys {
		keysInterface = append(keysInterface, v)
	}
	_, err = Command(nodeName, "DEL", keysInterface...)
	return
}

// exists command
func Exists(nodeName string, key string) (exist bool, err error) {
	exist, err = redis.Bool(Command(nodeName, "EXISTS", key))
	return
}

// expire command
func Expire(nodeName string, key string, seconds int) (err error) {
	_, err = Command(nodeName, "EXPIRE", key, seconds)
	return
}

// expireat command
func ExpireAt(nodeName string, key string, timestamp int64) (err error) {
	_, err = Command(nodeName, "EXPIREAT", key, timestamp)
	return
}

// keys command
func Keys(nodeName string, pattern string) (keys []string, err error) {
	keys, err = redis.Strings(Command(nodeName, "KEYS", pattern))
	return
}

// persist command
func Persist(nodeName, key string) (err error) {
	if _, err := Command(nodeName, "PERSIST", key); err != nil {
		return err
	}
	return nil
}

// ttl command
func TTL(nodeName, key string) (ttl int64, err error) {
	ttl, err = redis.Int64(Command(nodeName, "TTL", key))
	return
}

// setex command
func Setex(nodeName, key string, second int, value interface{}) (err error) {
	_, err = Command(nodeName, "SETEX", key, second, value)
	return
}

// setnx command
func Setnx(nodeName, key string, value interface{}) (err error) {
	_, err = Command(nodeName, "SETNX", key, value)
	return
}

// mset command
func MSet(nodeName string, valueMap map[string]interface{}) (err error) {
	values := make([]interface{}, 0, len(valueMap)*2)
	for k, v := range valueMap {
		values = append(values, k, v)
	}
	_, err = Command(nodeName, "MSET", values...)

	return
}

// mget command
func MGet(nodeName string, keys ...string) (values []interface{}, err error) {
	if len(keys) == 0 {
		return
	}
	keyInterface := make([]interface{}, 0, len(keys))
	for _, v := range keys {
		keyInterface = append(keyInterface, v)
	}
	values, err = redis.Values(Command(nodeName, "MGET", keyInterface...))
	return
}

// decr command
func Decr(nodeName, key string) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(nodeName, "DECR", key))
	return
}

// decrby command
func DecrBy(nodeName, key string, decrNum int) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(nodeName, "DECRBY", key, decrNum))
	return
}

// incr command
func Incr(nodeName, key string) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(nodeName, "INCR", key))
	return
}

// incrby command
func IncrBy(nodeName, key string, decrNum int) (newValue int64, err error) {
	newValue, err = redis.Int64(Command(nodeName, "INCRBY", key, decrNum))
	return
}

// getset command
func GetSet(nodeName, key string, value interface{}) (curValue interface{}, err error) {
	curValue, err = Command(nodeName, "GETSET", key, value)
	return
}

// hdel command
func HDel(nodeName string, key string, fields ...string) (err error) {
	if len(fields) == 0 {
		return
	}
	fieldsInterface := make([]interface{}, 0, len(fields)+1)
	fieldsInterface = append(fieldsInterface, key)
	for _, v := range fields {
		fieldsInterface = append(fieldsInterface, v)
	}
	_, err = Command(nodeName, "HDEL", fieldsInterface...)
	return
}

// hexist command
func HExists(nodeName, key, field string) (exist bool, err error) {
	exist, err = redis.Bool(Command(nodeName, "HEXISTS", key, field))
	return
}

// hset command
func HSet(nodeName, key, field string, value interface{}) (err error) {
	_, err = Command(nodeName, "HSET", key, field, value)
	return
}

// hmset command
func HMset(nodeName, key string, values interface{}) (err error) {
	args := redis.Args{}.Add(key)
	args = args.AddFlat(values)
	_, err = Command(nodeName, "HMSET", args...)
	return
}

// hget command
func HGet(nodeName, key, field string) (value interface{}, err error) {
	value, err = Command(nodeName, "HGET", key, field)
	return
}

// hgetall command
func HGetAll(nodeName, key string) (values []interface{}, err error) {
	values, err = redis.Values(Command(nodeName, "HGETALL", key))
	return
}

// hkeys command
func HKeys(nodeName, key string) (values []string, err error) {
	values, err = redis.Strings(Command(nodeName, "HKEYS", key))
	return
}

// hvals command
func HVals(nodeName, key string) (values []interface{}, err error) {
	values, err = redis.Values(Command(nodeName, "HVALS", key))
	return
}

// hlen command
func HLen(nodeName, key string) (length int64, err error) {
	length, err = redis.Int64(Command(nodeName, "HLEN", key))
	return
}

// hmget command
func HMget(nodeName, key string, fields ...string) (values interface{}, err error) {
	if len(fields) == 0 {
		return
	}
	fieldsInterface := make([]interface{}, 0, len(fields)+1)
	fieldsInterface = append(fieldsInterface, key)
	for _, v := range fields {
		fieldsInterface = append(fieldsInterface, v)
	}
	values, err = Command(nodeName, "HMGET", fieldsInterface...)
	return
}
