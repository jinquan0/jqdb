/*
	ln -s /root/jqdb/jqredis /usr/lib/golang/src/jqredis
	go mod init jqredis
	go mod tidy
	go build cmd/jqredis/list_main.go
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

	q := &jq.ST_RedisList {
		List: 	"mylist0",
		Data:   "abc12345678",
	}

	r:=jq.RedisConn(*e)
	jq.RPush(r, q)
	len:=jq.LLen(r, q); fmt.Printf("list length: %d\n", len)
	str:=jq.LPop(r, q); fmt.Printf("list reply data: %s\n", str)
	jq.RedisDisConn(r)
}
