# redis-go

A simple redis client written by Golang.

## install

```bash
go get github.com/scofieldpeng/redis-go
```

## Quick start

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/scofieldpeng/config"
	"github.com/scofieldpeng/redis-go"
	"github.com/gomodule/redigo/redis"
)

func main(){
    	
	// we use the ini file as the config format
	// config package document, see https://github.com/scofieldpeng/config pls
	// app.ini file
	// 
    // # default is the node name
    // [redis_node_default]
    // # scheme is the connection struct of redis
    // # format: redis://username:password@host:port/dbname
    // scheme=redis://user:secret@localhost:6379/db
    // # if this node have node, add node name to there, multiple nodes joined by comma
    // slaves=slave1
    // 
    // # node slave1
    // [redis_node_slave1]
    // scheme=redis://user:secret@localhost:6379/db

    // # node slave2
    // [redis_node_slave2]
    // scheme=redis://user:secret@localhost:6379/db
    nodesConfig := config.Data("app")
    goredis.Init(goredis.Config{
    	MaxActive:10,
    	MaxIdle:5,
    	Timeout:time.Second * 5,
        Wait:true,	
    },nodesConfig)
   
    _,err := goredis.Command("default","SET","name","scofield")
    if err == nil {
   	     // if u want to transfer result to string, use the redis.String() to get the result
   	     // more transfer helper, u can see the doc of github.com/gomodule/redigo/redis
   	     name,_:= redis.String(goredis.Command("default","GET","name"))
   	     fmt.Println(name)
    }
   
    // or USE the helper
    // equivalent to goredis.Command("default","SET","name","scofield")
    goredis.Set("default","name","scofield")
    // equivalent to goredis.Command("default","GET","name")
    goredis.Get("default","name")
     
    // more helper see https://github.com/scofieldpeng/redis-go/helper.go
}
```

## More documention

the redis package is based on [https://github.com/gomodule/redigo](https://github.com/gomodule/redigo), so you can view the detail from there:-)

## Licence
 
MIT Licence

## Thanks

1. [https://github.com/gomodule/redigo](https://github.com/gomodule/redigo)