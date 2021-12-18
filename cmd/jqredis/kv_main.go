/*
	ln -s /root/jqdb/jqredis /usr/lib/golang/src/jqredis
	go mod init jqredis
	go mod tidy
	go build cmd/jqredis/kv_main.go
*/


package main
 
import (
  jq "jqredis"
  "fmt"
)

func main() {
	e := &jq.ST_Redis_Endpoint {
		Endpoint: "172.24.16.13:26379",
		Pass:     "passw0rd",
		Db: 0,
	}

	kv := &jq.ST_RedisKv {
		Key: 	"myname",
		Val:   "jinquan",
	}

	r:=jq.RedisConn(*e)
	jq.KvSet(r, kv)
	v:=jq.KvGet(r, "myname"); fmt.Printf("Key/Value Pair reply data: %s\n", v)
	jq.RedisDisConn(r)
}
