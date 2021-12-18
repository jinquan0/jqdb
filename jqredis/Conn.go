package jqredis

import (
    "fmt"
    "github.com/go-redis/redis"
)

type ST_Redis_Endpoint struct {
  Endpoint  string
  Pass      string
  Db        int
}

func RedisConn(conn ST_Redis_Endpoint) *redis.Client {
    var RedisClnt  *redis.Client
	RedisClnt = redis.NewClient(&redis.Options{
        Addr:       conn.Endpoint,
        Password: 	conn.Pass,
        DB:       	conn.Db,
    })

    //ping
    pong, err := RedisClnt.Ping().Result()
    if err != nil {
        fmt.Println("Redis/> ping error", err.Error())
        return nil
    }
    fmt.Println("Redis/> ping result:", pong)
    return RedisClnt
}

func RedisDisConn(r *redis.Client) {
	r.Close()
}