/*
	ln -s /root/jqdb/jqmysql /usr/lib/golang/src/jqmysql
	go mod init jqmysql
	go mod tidy
	go build cmd/jqmysql/update_main.go
*/


package main
 
import (
  jq "jqmysql"
)


/*
@@@ 结构体元素数量 必须与 DELETE语句中的条件变量保持一致
@@@ 排序也必须保持一致
*/
type ST_MyFields_1 struct{
    Verbose_name   string  `json:"verbose_name"`
    Id         int    `json:"id"`
    Product_key    string  `json:"product_key"`
}

func main() {
	e := &jq.ST_MySQL_Endpoint {
	  Host:  "172.24.16.215",
	  Port:  "32294",
	  User:  "dba",
	  Pass:  "P@ssw0rd",
	  Sslca: "null",  // 不使用 SSL 连接
	  Db:    "test01",
	}

	d:=&ST_MyFields_1{
		Verbose_name: "jinquan",
		Id: 10,
		Product_key: "jjj",
	}

	db := jq.MyConn(e)
	jq.MyDelete( db, 
			"UPDATE bledev SET `verbose_name`=? WHERE `id`=? AND `product_key`=?",	*d ) 
	jq.MyDisconn(db)

}