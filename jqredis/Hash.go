/*
 https://www.jianshu.com/p/c9fa5540534d
*/
package jqredis

import (
    "fmt"
    "reflect"
    "time"
    "github.com/go-redis/redis"
)

func StructToMap(in interface{}, tagName string) (map[string]interface{}, error){
    out := make(map[string]interface{})

    v := reflect.ValueOf(in)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {  // 非结构体返回错误提示
        return nil, fmt.Errorf("Redis/> ToMap only accepts struct or struct pointer; got %T", v)
    }

    t := v.Type()
    // 遍历结构体字段
    // 指定tagName值为map中key;字段值为map中value
    for i := 0; i < v.NumField(); i++ {
        fi := t.Field(i)
        if tagValue := fi.Tag.Get(tagName); tagValue != "" {
            out[tagValue] = v.Field(i).Interface()
        }
    }
    return out, nil
}

func HashKeyExsit(endpoint ST_Redis_Endpoint, key string)  {
	r := RedisConn(endpoint)
	exsit := r.Exists(key)
	RedisDisConn(r)
	fmt.Println("Redis/> Hash key[%s] ",key,exsit)
}

func HashSet(endpoint ST_Redis_Endpoint, key string, fields map[string]interface{}) string {
    r := RedisConn(endpoint)
    val, err := r.HMSet(key, fields).Result()
    if err != nil {
        fmt.Println("Redis/> HMSet Error:", err)
    }
    RedisDisConn(r)
    return val
}

func HashSet_v2(endpoint ST_Redis_Endpoint, key string, fields map[string]interface{}, ttl time.Duration) string {
    r := RedisConn(endpoint)
    val, err := r.HMSet(key, fields).Result()
    if err != nil {
        fmt.Println("Redis/> HMSet Error:", err)
    }else{
	ttlcmd := r.Expire(key, ttl)
	    fmt.Printf("Redis/> hash key[%s] set TTL %v , result:%v\n", key, ttl, ttlcmd.Result())
    }
    RedisDisConn(r)
    return val
}

// 通过key 获取hash所有元素值
// 测试Hash 写入
// HMSET hk1 topic "aws-aurora-memory" key "nxprbdyxy-reader1" value "32"
func HashGetFields(r *redis.Client, key string, fields []string ) map[string]interface{} {
    resMap := make(map[string]interface{})
    for _, field := range fields {
        var result interface{}
        val, err := r.HGet(key, fmt.Sprintf("%s", field)).Result()
        if err == redis.Nil {
            fmt.Printf("Redis/> Key Doesn't Exists: %v\n", field)
            resMap[field] = result
        }else if err != nil {
            fmt.Printf("Redis/> HMGet Error: %v\n", err)
            resMap[field] = result
        }
        if val != "" {
            resMap[field] = val
        }else {
            resMap[field] = result
        }
    }
    return resMap
}

// 通过key 获取hash所有元素值
// 测试Hash 写入
// HMSET hk1 topic "aws-aurora-memory" key "nxprbdyxy-reader1" value "32"
func HashGetFields_v2(endpoint ST_Redis_Endpoint, key string, fields []string ) map[string]interface{} {
    r := RedisConn(endpoint)
    if r == nil {
	return nil
    }
    resMap := make(map[string]interface{})
    for _, field := range fields {
        var result interface{}
        val, err := r.HGet(key, fmt.Sprintf("%s", field)).Result()
        if err == redis.Nil {
            fmt.Printf("Redis/> Key Doesn't Exists: %v\n", field)
            resMap[field] = result
	    return nil
        }else if err != nil {
            fmt.Printf("Redis/> HMGet Error: %v\n", err)
            resMap[field] = result
        }
        if val != "" {
            resMap[field] = val
        }else {
            resMap[field] = result
        }
    }
    RedisDisConn(r)
    return resMap
}

func HashesGet(endpoint ST_Redis_Endpoint, keys []string, fields []string ) []map[string]interface{} {
    hashes := make([]map[string]interface{}, 1)
    hashes_cnt := len(keys)
    r := RedisConn(endpoint)
    for i:=0; i < hashes_cnt; i++ {
        if i == 0 {
            hashes[0] = HashGetFields(r, keys[0], fields) 
        } else if i > 0 {
            // first argument to append must be slice
            hashes = append(hashes, HashGetFields(r, keys[i], fields) )
        }
    }
    RedisDisConn(r)
    return hashes
}



func HashesPrint(hashes []map[string]interface{}) {
	for i:=0; i < len(hashes); i++ {
		fmt.Println(hashes[i])
	}
}
