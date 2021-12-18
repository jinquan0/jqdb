/*
	ln -s /root/jqdb/jqmysql /usr/lib/golang/src/jqmysql
	go mod init jqmysql
	go mod tidy
	go build cmd/jqmysql/delete_main.go
*/


package main
 
import (
  jq "jqmysql"
)


/*
@@@ 结构体元素数量 必须与 DELETE语句中的条件变量保持一致
*/
type ST_MyFields_1 struct{
    //Id         int    `json:"id"`
    Product_key    string  `json:"product_key"`
    Verbose_name   string  `json:"verbose_name"`
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
		Product_key: "jinquan",
		Verbose_name: "k8s",
	}

	db := jq.MyConn(e)
	jq.MyDelete( db, 
			"DELETE FROM bledev WHERE `product_key`=? AND `verbose_name`=?",	*d ) 
	jq.MyDisconn(db)

}