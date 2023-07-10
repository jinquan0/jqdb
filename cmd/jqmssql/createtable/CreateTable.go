package main

import (
	"context"
	"log"
	"reflect"
    "flag"
    _ "github.com/denisenkom/go-mssqldb"
    jq "github.com/jinquan0/jqdb/jqmssql"
    //jq "gitee.com/jinquan711/jdb/jqmssql"
)

var (
    debug         = flag.Bool("debug", false, "enable debugging")
    server        = flag.String("server", "172.24.22.4", "SQL Server Host IP")
    port     *int = flag.Int("port", 1433, "SQL Server connect port")
    user          = flag.String("user", "testuser0", "connect user")
    password      = flag.String("password", "FDNvQ8#g", "connect password")
    //user          = flag.String("user", "cn\\infra01testuser", "connect user")
    database      = flag.String("database", "jq01test", "database name")
)

func init() {
	flag.Parse()
}

func CreateTable() {
    conn := jq.MssqlConn(*server, *port, *user, *password, *database)
    //conn := jq.MssqlConnWithMSA(*server, *port, *user, *database)
	_, err := conn.Exec("create table test (f int)")
	defer conn.Exec("drop table test")
	if err != nil {
		log.Fatal("create table failed with error", err)
	}
	jq.MssqlDisconn(conn)
}


func Transaction() {
    conn := jq.MssqlConn(*server, *port, *user, *password, *database)
    //conn := jq.MssqlConnWithMSA(*server, *port, *user, *database)
	tx1, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatal("BeginTx failed with error", err)
	}
	tx2, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatal("BeginTx failed with error", err)
	}
	_, err = tx1.Exec("insert into test (f) values (1)")
	if err != nil {
		log.Fatal("insert failed with error", err)
	}
	_, err = tx2.Exec("insert into test (f) values (2)")
	if err != nil {
		log.Fatal("insert failed with error", err)
	}
	tx1.Rollback()
	tx2.Commit()

	rows, err := conn.Query("select f from test")
	if err != nil {
		log.Fatal("select failed with error", err)
	}
	defer rows.Close()
	values := []int64{}
	for rows.Next() {
		var val int64
		err = rows.Scan(&val)
		if err != nil {
			log.Fatal("scan failed with error", err)
		}
		values = append(values, val)
	}
	if !reflect.DeepEqual(values, []int64{2}) {
		log.Fatal("Values is expected to be [1] but it is %v", values)
	}
	jq.MssqlDisconn(conn)
}

func main() {
	CreateTable()
	Transaction()
}
