package main

import (
    "flag"
    "fmt"
    _ "github.com/denisenkom/go-mssqldb"
    jq "github.com/jinquan0/jqdb/jqmssql"
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
    Id int64  `json:"id"`
    Fld1 int64 `json:"fld1"`
    Fld2 string  `json:"fld2"`
}

func AnyarrayAlloc(sz int) []interface{} {
    any_array := make([]interface{}, sz, sz)
    for i:=0; i < sz; i++ {
        any_array[i]=new(MsQueryReply)
    }
    return any_array
}


var DBC_PARA jq.DBconn

func init() {
    flag.Parse()
    DBC_PARA.Server = *server
    DBC_PARA.Port = *port 
    DBC_PARA.User = *user 
    DBC_PARA.Pass = *password
    DBC_PARA.Database = *database
}

func main() {
    //db:=jq.MssqlConn(*server, *port, *user, *password)
    db := jq.MssqlConn(&DBC_PARA)

    //SELECT [Id],[F1],[F2],[F3],[F4],[F5] FROM [DB0].[dbo].[Table_1] WHERE Id=3 OR Id=5
    q:=&jq.ST_MsCommonQuery {
        Sql: `select id,fld1,fld2 from table1 where id>=101 and id<=111`,
        AnyArray: AnyarrayAlloc(16),
    }
    _,MAPS,numrow:=jq.MsCommonSelect2Maps(*q, db)
    jq.MssqlDisconn(db)

    //fmt.Println(numrow)
    for i:=0; i < numrow; i++ {
        fmt.Println(MAPS[i])
    }
}
