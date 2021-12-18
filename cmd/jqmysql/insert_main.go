/*
	ln -s /root/jqdb/jqmysql /usr/lib/golang/src/jqmysql
	go mod init jqmysql
	go mod tidy
	go build cmd/jqmysql/insert_main.go
*/


package main
 
import (
  jq "jqmysql"
  //"fmt"
)


/*
表结构Field name 与结构体json 标签一致(字符大小写)
+--------------+--------------+------+-----+---------+-------+
| Field        | Type         | Null | Key | Default | Extra |
+--------------+--------------+------+-----+---------+-------+
| id           | int(11)      | NO   | PRI | NULL    |       |
| product_key  | varchar(255) | YES  |     | NULL    |       |
| verbose_name | varchar(255) | YES  |     | NULL    |       |
+--------------+--------------+------+-----+---------+-------+
*/
type ST_MyFields_1 struct{
    Id         int    `json:"id"`
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
		Id: 11,
		Product_key: "jinquan",
		Verbose_name: "k8s",
	}

	db := jq.MyConn(e)
	jq.MyInsert( db, 
			"INSERT IGNORE INTO bledev(`id`, `product_key`, `verbose_name`) VALUES (?,?,?)",	*d ) 
	jq.MyDisconn(db)

}