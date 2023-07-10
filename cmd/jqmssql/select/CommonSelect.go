package main

import (
    "flag"
    "fmt"
    _ "github.com/denisenkom/go-mssqldb"
    jq "github.com/jinquan0/jqdb/jqmssql"
    //jq "gitee.com/jinquan711/jdb/jqmssql"
)

var (
    debug         = flag.Bool("debug", false, "enable debugging")
    server        = flag.String("server", "172.24.22.4", "the database server")
    port     *int = flag.Int("port", 1433, "the database port")
    user          = flag.String("user", "testuser0", "the database user")
    password      = flag.String("password", "FDNvQ8#g", "the database password")
    database      = flag.String("database", "jq01test", "database name")
)

type MsQueryReply struct{
    Id int64  `json:"Id"`
    F1 string `json:"F1"`
    F2 string  `json:"F2"`
    F3 string  `json:"F3"`
    F4 float64  `json:"F4"`
    F5 float64  `json:"F5"`
}
func AnyarrayAlloc(sz int) []interface{} {
    any_array := make([]interface{}, sz, sz)
    for i:=0; i < sz; i++ {
        any_array[i]=new(MsQueryReply)
    }
    return any_array
}

func main() {
    flag.Parse()
    db:=jq.MssqlConn(*server, *port, *user, *password)

    q:=&jq.ST_MsCommonQuery {
        Sql: `SELECT [Id],[F1],[F2],[F3],[F4],[F5] FROM [DB0].[dbo].[Table_1] WHERE Id=3 OR Id=5`,
        AnyArray: AnyarrayAlloc(16),
    }
    _,MAPS,numrow:=jq.MsCommonSelect2Maps(*q, db)
    jq.MssqlDisconn(db)

    fmt.Println(numrow)
    fmt.Println(MAPS)

}
