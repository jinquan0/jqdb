package jqredis

import (
    "fmt"
    "time"
    "github.com/tidwall/gjson"
    "github.com/go-redis/redis"
)

const (
    REDIS_LIST_PUSH_OK = 0
    REDIS_LIST_POP_FAIL = 1
    REDIS_LIST_LEN_FAIL = 2
)

type (
  ST_RedisList struct {
    List  string `json:"List" validate:"required"`
    Data  string `json:"Data" validate:"required"`
  }
)

// Redis list : Queue
//      		-------------------------------------
//      <- LPop	| *	| *	| *	| *	| *	| *	| *	| *	| *	|  <- RPush
//      		-------------------------------------
func RPush(r *redis.Client, l *ST_RedisList) {
    r.RPush(l.List, l.Data).Err()
}

func LPop(r *redis.Client, l *ST_RedisList) string {
	first, _ := r.LPop(l.List).Result()
	return first
}

func BLPop(r *redis.Client, l *ST_RedisList, timeout_s time.Duration) []string {
	first, _ := r.BLPop(time.Second*timeout_s, l.List).Result()
	return first
}

func LLen(r *redis.Client, l *ST_RedisList) int64 {
	listLen, _ := r.LLen(l.List).Result()
	return listLen
}

// RPush with DeDuplication
// 具备list key消重功能
func RPushDeDup(r *redis.Client, l *ST_RedisList, key_DeDup string) {
	val_DeDup:=gjson.Get(l.Data, key_DeDup)
	list,err:=r.LRange(l.List, 0, LLen(r, l)).Result()
	if err!=nil {
		fmt.Println(err)
		return
	}
	var match_flag bool; match_flag=false
	for _,v := range list {
		v_DeDup:=gjson.Get(v, key_DeDup)
		//fmt.Printf("%s -- %s\n", val_DeDup.String(), v_DeDup.String())
		if val_DeDup.String() == v_DeDup.String() {
			match_flag=true
		}
	}
	if match_flag == false {
		RPush(r, l)
		fmt.Println("Redis/> RPush, DeDup-Key[%s]", key_DeDup)
	}

}