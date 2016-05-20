# redis-go

A simple redis client written by Golang.

## install

```bash
go get github.com/scofieldpeng/redis-go
```

## Usage

First, create ini File in your `$Project_Directory/config/redis.ini`, write these config the `redis.ini`

```ini
[config]
# max idle connection in pools
maxIdle=5
# idle timeout(seconds) in pools
idleTimeout=60
# connect && read && write timeout(seconds)
timeout=10

# redis server nodes
[nodes]
# default is node name
default=127.0.0.1:6379?password=123456
slave1=127.0.0.1:6378?password=123456
```

Second, Load ini file

```go
redisConfigs := config.Config("redis") // import github.com/scofieldpeng/config-go,this config docuement see the github :-)
```

Now, Init Config

```go
redis.Init(redisConfigs)
```

then you can use

```go
conn := redis.Pool().Get()
defer conn.Close()

if _,err := conn.Do("SET","test",123);err != nil {
    log.Println(err.Error())
}
```

## More documention

the redis package is based on [https://github.com/garyburd/redigo](https://github.com/garyburd/redigo), so you can view the detail from there:-)

## Licence
 
MIT Licence

## Thanks

1. [https://github.com/garyburd/redigo](https://github.com/garyburd/redigo)