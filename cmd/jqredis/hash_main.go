/*
	ln -s /root/jqdb/jqredis /usr/lib/golang/src/jqredis
	go mod init jqredis
	go mod tidy
	go build cmd/jqredis/hash_main.go
*/


package main
 
import (
  jq "jqredis"
)

type ST_MyHash struct {
  Topic  string `json:"topic"`
  Key    string `json:"key"`
  Value    int  `json:"value"`
}

func main() {
	e := &jq.ST_Redis_Endpoint {
		Endpoint: "172.24.16.13:26379",
		Pass:     "passw0rd",
		Db: 0,
	}

	h1 := ST_MyHash{Topic: "PCIe-card", Key: "Intel", Value: 8}
	m, _ := jq.StructToMap(&h1, "json")
	jq.HashSet(*e, "hk5", m )

	k := []string {"hk1","hk2","hk3", "hk4", "hk5"}
	f := []string {"topic","key","value"}
	h := jq.HashesGet(*e, k, f)
	jq.HashesPrint(h) 

}
