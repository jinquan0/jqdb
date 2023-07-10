package main

import (
	//"bytes"
	"context"
	//"database/sql"
	//"database/sql/driver"
	//"fmt"
	//"io"
	"log"
	//"math"
	//"net"
	"reflect"
	//"strings"
	//"sync"
	//"testing"
	//"time"

	// "github.com/denisenkom/go-mssqldb/msdsn"

    "flag"
    //"fmt"
    _ "github.com/denisenkom/go-mssqldb"
    jq "github.com/jinquan0/jqdb/jqmssql"

)

var (
    debug         = flag.Bool("debug", false, "enable debugging")
    password      = flag.String("password", "n9TLs7wp", "the database password")
    port     *int = flag.Int("port", 1433, "the database port")
    server        = flag.String("server", "172.24.22.1", "the database server")
    user          = flag.String("user", "sa", "the database user")
)

func init() {
	flag.Parse()
}

func CreateTable() {
    conn := jq.MssqlConn(*server, *port, *user, *password)

	_, err := conn.Exec("create table test (f int)")
	defer conn.Exec("drop table test")
	if err != nil {
		log.Fatal("create table failed with error", err)
	}
	jq.MssqlDisconn(conn)
}


func Transaction() {
	//conn, logger := open(t)
	//defer conn.Close()
	//defer logger.StopLogging()

    conn := jq.MssqlConn(*server, *port, *user, *password)

	_, err := conn.Exec("create table test (f int)")
	defer conn.Exec("drop table test")
	if err != nil {
		log.Fatal("create table failed with error", err)
	}

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
