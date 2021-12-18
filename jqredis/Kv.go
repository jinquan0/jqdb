package jqredis

import (
    "fmt"
    "github.com/go-redis/redis"
)

type (
  ST_RedisKv struct {
    Key  string `json:"Key" validate:"required"`
    Val  string `json:"Val" validate:"required"`
  }
)

// Redis key-value pair
func KvSet(r *redis.Client, kv *ST_RedisKv) {
    err := r.Set(kv.Key, kv.Val, 0).Err()
    if err != nil {
        fmt.Println(err)
    }
}

func KvGet(r *redis.Client, k string) string {
    val,err := r.Get(k).Result()
    if err != nil {
        fmt.Println(err)
        return ""
    }
    return val 
}