/*
	ln -s /root/jqdb/jqmysql /usr/lib/golang/src/jqmysql
	go mod init jqmysql
	go mod tidy
	go build cmd/jqmysql/common_select_main.go
*/


package main
 
import (
  jq "jqmysql"
  "fmt"
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

func AnyarrayAlloc(sz int) []interface{} {
    any_array := make([]interface{}, sz, sz)
    for i:=0; i < sz; i++ { 
        any_array[i]=new(ST_MyFields_1) 
    }
    return any_array
}

func ST_MyFields_1_Print(array []interface{}, num int) {
    for i:=0; i < num; i++ {
        fmt.Printf("id: %d\t product_key: %s\t verbose_name:%s\n", 
            jq.ST_GetValueByKey(array[i], "Id"), 
            jq.ST_GetValueByKey(array[i], "Product_key"),
            jq.ST_GetValueByKey(array[i], "Verbose_name") )
    }
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

	db := jq.MyConn(e)

	q:=&jq.ST_MyCommonQuery {
		Sql: "select id,product_key,verbose_name from bledev;",
		AnyArray: AnyarrayAlloc(16),
	}
	_,numrow:=jq.MyCommonSelect(*q, db)

	jq.MyDisconn(db)

	ST_MyFields_1_Print(q.AnyArray, numrow)

}